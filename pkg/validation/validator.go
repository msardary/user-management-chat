package validation

import (
	"errors"

	locale_en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	validator_en "github.com/go-playground/validator/v10/translations/en"
)


var (
	Validate *validator.Validate
	Trans ut.Translator
)


func InitValidator() error{
	Validate = validator.New()

	enLocale := locale_en.New()
	uni := ut.New(enLocale, enLocale)

	trans, found := uni.GetTranslator("en")
	if !found {
		return errors.New("translator not found")
	}

	Trans = trans

	return validator_en.RegisterDefaultTranslations(Validate, Trans)
}