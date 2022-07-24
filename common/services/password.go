package services

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/v4lproik/simple-blockchain-quickstart/common/utils"

	"golang.org/x/crypto/argon2"
)

// inspired from: https://gist.github.com/alexedwards/34277fae0f48abe36822b375f0f6a621
var (
	ErrInvalidHash         = errors.New("invalid format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

type PasswordService struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewPasswordService(memory uint32, iterations uint32, parallelism uint8, saltLength uint32, keyLength uint32) PasswordService {
	return PasswordService{memory: memory, iterations: iterations, parallelism: parallelism, saltLength: saltLength, keyLength: keyLength}
}

func NewDefaultPasswordService() PasswordService {
	return PasswordService{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

func (c PasswordService) GenerateHash(password string) (hash string, err error) {
	salt, err := utils.GenerateRandomBytes(c.saltLength)
	if err != nil {
		return "", err
	}

	hashx := argon2.IDKey([]byte(password), salt, c.iterations, c.memory, c.parallelism, c.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hashx)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, c.memory, c.iterations, c.parallelism, b64Salt, b64Hash), nil
}

func (c PasswordService) ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password hash.
	salt, hash, err := c.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, c.iterations, c.memory, c.parallelism, c.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	return subtle.ConstantTimeCompare(hash, otherHash) == 1, nil
}

func (p PasswordService) decodeHash(encodedHash string) ([]byte, []byte, error) {
	var err error
	var version int
	var salt []byte
	var hash []byte

	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, ErrInvalidHash
	}

	if _, err = fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, ErrIncompatibleVersion
	}

	if _, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism); err != nil {
		return nil, nil, err
	}

	if salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4]); err != nil {
		return nil, nil, err
	}

	if hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5]); err != nil {
		return nil, nil, err
	}

	return salt, hash, nil
}
