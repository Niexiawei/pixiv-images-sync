package validator

import (
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	chiese "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validatorTrans  ut.Translator
	validatorEngine = binding.Validator.Engine().(*validator.Validate)
)

func InitValidatorTrans() {
	zhCn := zh.New()
	uni := ut.New(zhCn, zhCn)
	validatorTrans, _ = uni.GetTranslator("zh")
	_ = chiese.RegisterDefaultTranslations(validatorEngine, validatorTrans)
}

type RequestBindingError map[string]string

func FormatValidatorErrors(err error) *RequestBindingError {
	if err, ok := err.(validator.ValidationErrors); ok {
		errors := make(RequestBindingError)
		for _, value := range err {
			errors[value.Field()] = value.Translate(validatorTrans)
		}
		return &errors
	}
	return &RequestBindingError{
		"errors": "字段验证错误",
	}
}

func (r *RequestBindingError) All() map[string]string {
	return *r
}

func (r *RequestBindingError) First() string {
	errorValues := maputil.Values[string, string](*r)
	if len(errorValues) < 1 {
		return ""
	}
	return errorValues[0]
}
