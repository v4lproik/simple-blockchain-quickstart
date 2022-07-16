package services

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/thoas/go-funk"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
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

func checkVerifyingConf(conf VerifyingConf) error {
	if conf.jkmsRefreshCacheIntervalInMin == 0 {
		return errors.New("checkVerifyingConf: refresh cache interval in minute cannot be equal to 0")
	}

	if conf.jkmsRefreshCacheTimeoutInSec == 0 {
		return errors.New("checkVerifyingConf: refresh cache timeout in second cannot be equal to 0")
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

func checkSigningConf(conf SigningConf) error {
	algo := conf.algo
	if !funk.Contains(ALLOWED_ALGORITHMS, algo) {
		return errors.New("checkSigningConf: jwt algo " + algo + " is not allowed")
	}

	privateKeyPath := conf.privateKeyPath
	_, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("checkSigningConf: failed to read private key: %w", err)
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
	err := checkSigningConf(signingConf)
	if err != nil {
		return nil, fmt.Errorf("NewJwtService: invalid signing configuration: %w", err)
	}

	err = checkVerifyingConf(verifyingConf)
	if err != nil {
		return nil, fmt.Errorf("NewJwtService: invalid verifying configuration: %w", err)
	}

	// get the private key
	privateKeyPath := signingConf.privateKeyPath
	privateKey, _ := ioutil.ReadFile(privateKeyPath)

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, fmt.Errorf("NewJwtService: failed at parsing private key: %w", err)
	}

	return &JwtService{
		verifyingConf: verifyingConf,
		signingConf:   signingConf,
		privateKey:    key,
		jwksClient:    nil,
	}, nil
}

// SignToken Sign token with private key passed at the initialisation of the service
// including the payload passed as content parameter
func (j *JwtService) SignToken(content interface{}) (string, error) {
	var signedToken string
	signingConf := j.signingConf
	now := time.Now().UTC()

	payload, err := json.Marshal(content)
	if err != nil {
		return signedToken, fmt.Errorf("SignToken: failed at marshaling payload: %w", err)
	}

	claims := make(jwt.MapClaims)
	claims["dat"] = string(payload)                                                       // Our custom data.
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

	signedToken, err = token.SignedString(j.privateKey)
	if err != nil {
		return signedToken, fmt.Errorf("SignToken: failed to sign the token: %w", err)
	}

	return signedToken, nil
}

// VerifyToken check if token is valid according to public key info expsed by jwks
func (j *JwtService) VerifyToken(signedToken string) (jwt.Token, error) {
	//lazy initialisation singleton
	//very useful in the case of the client needs to be initiated after the endpoint is exposed
	if j.jwksClient == nil {
		err := j.getJwksClient()
		if err != nil {
			return jwt.Token{}, fmt.Errorf("VerifyToken: failed to get jwks client: %w", err)
		}
	}

	//parse token and validate it (eg is expired?)
	token, err := jwt.Parse(signedToken, j.jwksClient.Keyfunc)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("VerifyToken: failed to parse token: %w", err)
	}
	return *token, nil
}

func (j *JwtService) PrivateKeyPath() string {
	return j.signingConf.privateKeyPath
}

func (j *JwtService) PrivateKeyId() string {
	return j.signingConf.privateKeyId
}

//initiate jwks client
func (j *JwtService) getJwksClient() error {
	var err error
	// initiate the jwks client
	options := keyfunc.Options{
		RefreshInterval:  time.Minute * time.Duration(j.verifyingConf.jkmsRefreshCacheIntervalInMin),
		RefreshTimeout:   time.Second * time.Duration(j.verifyingConf.jkmsRefreshCacheTimeoutInSec),
		RefreshRateLimit: time.Second * time.Duration(j.verifyingConf.jkmsRefreshCacheTimeoutInSec),
		RefreshErrorHandler: func(err error) {
			Logger.Errorf("jwks client cannot reach the jwks service: %s", err)
		},
	}

	// create the JWKS from the resource at the given URL.
	j.jwksClient, err = keyfunc.Get(j.verifyingConf.jwksUrl, options)
	if err != nil {
		return fmt.Errorf("getJwksClient: failed to initiate jwks client: %w", err)
	}

	return nil
}
