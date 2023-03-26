package validator

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validator  *validator.Validate
	translator ut.Translator
}

func NewValidator(translator ut.Translator) *Validator {
	return &Validator{
		validator:  validator.New(),
		translator: translator,
	}
}

func (v *Validator) GetTranslator() ut.Translator {
	return v.translator
}

func (v *Validator) Struct(s interface{}) error {
	return v.validator.Struct(s)
}

func NewEnglishTranslator() ut.Translator {
	en := en.New()
	uni := ut.New(en, en)
	translator, _ := uni.GetTranslator("en")

	return translator
}

func (v *Validator) ConfigureValidator() {
	en_translations.RegisterDefaultTranslations(v.validator, v.translator)
	v.RegisterTranslations()
	v.RegisterTagNameFunc(lowerCaseTagNameFunction)
	v.RegisterCustomTagValidations()
}

func (v *Validator) RegisterCustomTagValidations() {
	v.validator.RegisterValidation("receiptUrl", receiptUrlValidation)
}

func (v *Validator) RegisterTranslations() {
	// v.RegisterTranslation("required", "{0} is required.")
	// v.RegisterTranslation("email", "{0} must be a valid email.")
	// v.RegisterTranslation("max", "{0} has maximum limit of {1}")
	// v.RegisterTranslation("min", "{0} should have at least {1} character(s).")

	_ = v.validator.RegisterTranslation("required", v.translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.validator.RegisterTranslation("email", v.translator, func(ut ut.Translator) error {
		return ut.Add("email", "{0} should be a valid email.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = v.validator.RegisterTranslation("max", v.translator, func(ut ut.Translator) error {
		return ut.Add("max", "{0} has a limit of maximum {1} characters.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.Field(), fe.Param())
		return t
	})

	_ = v.validator.RegisterTranslation("min", v.translator, func(ut ut.Translator) error {
		return ut.Add("min", "{0} must have at least {1} characters.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field(), fe.Param())
		return t
	})

	_ = v.validator.RegisterTranslation("receiptUrl", v.translator, func(ut ut.Translator) error {
		return ut.Add("receiptUrl", "{0} is not valid.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("receiptUrl", fe.Field())
		return t
	})
}

func (v *Validator) RegisterTagNameFunc(f validator.TagNameFunc) {
	v.validator.RegisterTagNameFunc(f)
}

func (v *Validator) RegisterTranslation(tag string, message string, params ...string) {
	v.validator.RegisterTranslation(
		tag,
		v.translator,
		registerTranslationFunc(tag, message),
		translationFunc(tag, params...),
	)
}

func registerTranslationFunc(tag string, message string) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, message, true)
	}
}

func translationFunc(tag string, params ...string) validator.TranslationFunc {
	return func(ut ut.Translator, fe validator.FieldError) string {
		params = append([]string{fe.Field()}, params...)
		t, _ := ut.T(tag, params...)

		return t
	}
}

func lowerCaseTagNameFunction(field reflect.StructField) string {
	name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}

	return name
}

// Custom tag validation functions
func receiptUrlValidation(fl validator.FieldLevel) bool {
	const FISACLIZATION_SYSTEM_HOST = "suf.purs.gov.rs"
	url, err := url.Parse(fl.Field().String())
	if err != nil {
		return false
	}

	hostname := strings.TrimPrefix(url.Hostname(), "www.")

	return hostname == FISACLIZATION_SYSTEM_HOST
}
