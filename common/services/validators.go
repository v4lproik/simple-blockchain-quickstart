package services

import (
	"unicode"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
)

type ValidatorService struct{}

type Enum interface {
	IsValid() bool
}

func ValidateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(Enum)
	return value.IsValid()
}

func isValidPassword(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func ValidatePassword(fl validator.FieldLevel) bool {
	if pwd, ok := fl.Field().Interface().(string); ok {
		if !isValidPassword(pwd) {
			return false
		}
	}
	return true
}

func ValidateAccount(fl validator.FieldLevel) bool {
	if account, ok := fl.Field().Interface().(string); ok {
		if !common.IsHexAddress(account) {
			return false
		}
	}
	return true
}

func ValidateHash(fl validator.FieldLevel) bool {
	if hash, ok := fl.Field().Interface().(string); ok {
		// check blocks model for length
		if len(hash) == 64 {
			return true
		}
	}
	return false
}

func (v ValidatorService) AddValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("enum", ValidateEnum)
		if err != nil {
			Logger.Fatalf("AddValidators: custom validator enum cannot be registered: %w", err)
		}
		err = v.RegisterValidation("password", ValidatePassword)
		if err != nil {
			Logger.Fatalf("AddValidators: custom validator password cannot be registered: %w", err)
		}
		err = v.RegisterValidation("account", ValidateAccount)
		if err != nil {
			Logger.Fatalf("AddValidators: custom validator account cannot be registered: %w", err)
		}
		err = v.RegisterValidation("hash", ValidateHash)
		if err != nil {
			Logger.Fatalf("AddValidators: custom validator hash cannot be registered: %w", err)
		}
	}
}
