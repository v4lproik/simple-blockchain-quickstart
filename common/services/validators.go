package services

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	log "go.uber.org/zap"
)

type ValidatorService struct{}

type Enum interface {
	IsValid() bool
}

func ValidateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(Enum)
	return value.IsValid()
}

func (v ValidatorService) AddValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("enum", ValidateEnum)
		if err != nil {
			log.S().Fatalf("custom validators cannot be registered %v", err)
		}
	}
}
