package services

import (
	"crypto/rsa"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thoas/go-funk"
	log "go.uber.org/zap"
	"io/ioutil"
	"time"
)

var (
	ALLOWED_ALGORITHMS = []string{"HS256"}
)

type VerifyingConf struct {
	jwksUrl                        string
	jkmsRefreshCacheIntervalInMin  int
	jkmsRefreshCacheRateLimitInMin int
	jkmsRefreshCacheTimeoutInSec   int
}

func NewVerifyingConf(jwksUrl string, jkmsRefreshCacheIntervalInMin int, jkmsRefreshCacheRateLimitInMin int, jkmsRefreshCacheTimeoutInSec int) VerifyingConf {
	return VerifyingConf{jwksUrl: jwksUrl, jkmsRefreshCacheIntervalInMin: jkmsRefreshCacheIntervalInMin, jkmsRefreshCacheRateLimitInMin: jkmsRefreshCacheRateLimitInMin, jkmsRefreshCacheTimeoutInSec: jkmsRefreshCacheTimeoutInSec}
}

func checkVerifyingConfOrError(conf VerifyingConf) error {
	if conf.jkmsRefreshCacheIntervalInMin == 0 {
		return fmt.Errorf("refresh cache interval in minute cannot be equal to 0")
	}

	if conf.jkmsRefreshCacheTimeoutInSec == 0 {
		return fmt.Errorf("refresh cache timeout in second cannot be equal to 0")
	}
	return nil
}

type SigningConf struct {
	algo           string
	audience       string
	domain         string
	expiresInHours int
	issuer         string
	privateKeyPath string
	privateKeyId   string
}

func NewSigningConf(algo string, audience string, domain string, expiresInHours int, issuer string, privateKeyPath string, privateKeyId string) SigningConf {
	return SigningConf{algo: algo, audience: audience, domain: domain, expiresInHours: expiresInHours, issuer: issuer, privateKeyPath: privateKeyPath, privateKeyId: privateKeyId}
}

func checkSigningConfOrError(conf SigningConf) error {
	algo := conf.algo
	if !funk.Contains(ALLOWED_ALGORITHMS, algo) {
		return fmt.Errorf("jwt algo %v is not allowed", algo)
	}

	privateKeyPath := conf.privateKeyPath
	_, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("private key cannot be opened %v", err)
	}

	return nil
}

type JwtService struct {
	privateKey    *rsa.PrivateKey
	jwksClient    *keyfunc.JWKS
	signingConf   SigningConf
	verifyingConf VerifyingConf
}

func NewJwtService(verifyingConf VerifyingConf, signingConf SigningConf) (*JwtService, error) {
	// check the configurations for signing and verifying
	err := checkSigningConfOrError(signingConf)
	if err != nil {
		return nil, fmt.Errorf("error with signing conf %v", err)
	}

	err = checkVerifyingConfOrError(verifyingConf)
	if err != nil {
		return nil, fmt.Errorf("error with verifying conf %v", err)
	}

	// get the private key
	privateKeyPath := signingConf.privateKeyPath
	privateKey, _ := ioutil.ReadFile(privateKeyPath)

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, fmt.Errorf("private key cannot be parsed %v", err)
	}

	return &JwtService{
		verifyingConf: verifyingConf,
		signingConf:   signingConf,
		privateKey:    key,
		jwksClient:    nil,
	}, nil
}

func (j *JwtService) SignToken(content interface{}) (string, error) {
	signingConf := j.signingConf
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content                                                               // Our custom data.
	claims["exp"] = now.Add(time.Hour * time.Duration(signingConf.expiresInHours)).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()                                                            // The time at which the token was issued.
	claims["nbf"] = now.Unix()                                                            // The time before which the token must be disregarded.

	var signingMethod *jwt.SigningMethodRSA
	switch signingConf.algo {
	case "HS256":
		signingMethod = jwt.SigningMethodRS256
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	token.Header["kid"] = signingConf.privateKeyId

	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing the jwt token %v", err)
	}

	return signedToken, nil
}

func (j *JwtService) VerifyToken(signedToken string) (jwt.Token, error) {
	//lazy initialisation singleton
	//very useful in the case of the client needs to be initiated after the endpoint is exposed
	if j.jwksClient == nil {
		err := j.getJwksClient()
		if err != nil {
			return jwt.Token{}, fmt.Errorf("error verifying the jwt token %v", err)
		}
	}

	token, err := jwt.Parse(signedToken, j.jwksClient.Keyfunc)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("error verifying the jwt token %v", err)
	}
	return *token, nil
}

func (j *JwtService) PrivateKeyPath() string {
	return j.signingConf.privateKeyPath
}

func (j *JwtService) PrivateKeyId() string {
	return j.signingConf.privateKeyId
}

func (j *JwtService) getJwksClient() error {
	var err error
	// initiate the jwks client
	options := keyfunc.Options{
		RefreshInterval:  time.Minute * time.Duration(j.verifyingConf.jkmsRefreshCacheIntervalInMin),
		RefreshTimeout:   time.Second * time.Duration(j.verifyingConf.jkmsRefreshCacheTimeoutInSec),
		RefreshRateLimit: time.Second * time.Duration(j.verifyingConf.jkmsRefreshCacheTimeoutInSec),
		RefreshErrorHandler: func(err error) {
			log.S().Errorf("couldn't reach the jwks url %v", err.Error())
		},
	}

	// create the JWKS from the resource at the given URL.
	j.jwksClient, err = keyfunc.Get(j.verifyingConf.jwksUrl, options)
	if err != nil {
		return fmt.Errorf("error initiating jwks client %v", err)
	}

	return nil
}
