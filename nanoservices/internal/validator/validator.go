package validator

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"
	db "github.com/steffanturanjanin/receipt-manager/internal/database"
	"github.com/steffanturanjanin/receipt-manager/internal/errors"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

func NewValidator(translator ut.Translator) *Validator {
	return &Validator{
		Validator:  validator.New(),
		Translator: translator,
	}
}

func NewDefaultValidator() *Validator {
	validator := NewValidator(NewEnglishTranslator())
	validator.ConfigureValidator()

	return validator
}

func (v *Validator) ValidateEvent(event interface{}) error {
	err := v.Validate(event)
	if err == nil {
		return nil
	}

	validationErrors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		validationErrors[fieldName] = err.Translate(v.GetTranslator())
	}

	return &errors.HttpError{
		ErrBase: errors.ErrBase{Err: err, Message: "Request validation failed."},
		Code:    http.StatusBadRequest,
		Errors:  validationErrors,
	}
}

func (v *Validator) GetValidationErrors(s interface{}) map[string]string {
	err := v.Validate(s)
	if err == nil {
		return nil
	}

	validationErrors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		validationErrors[fieldName] = err.Translate(v.GetTranslator())
	}

	return validationErrors
}

func (v *Validator) GetTranslator() ut.Translator {
	return v.Translator
}

func (v *Validator) Validate(s interface{}) error {
	return v.Validator.Struct(s)
}

func NewEnglishTranslator() ut.Translator {
	en := en.New()
	uni := ut.New(en, en)
	translator, _ := uni.GetTranslator("sr")

	return translator
}

func (v *Validator) ConfigureValidator() {
	en_translations.RegisterDefaultTranslations(v.Validator, v.Translator)
	v.RegisterTranslations()
	v.RegisterTagNameFunc(lowerCaseTagNameFunction)
	v.RegisterCustomTagValidations()
}

func (v *Validator) RegisterCustomTagValidations() {
	v.Validator.RegisterValidation("receipt_url", receiptUrlValidation)
	v.Validator.RegisterValidation("url_query_params", urlQueryParamsValidation)
	v.Validator.RegisterValidation("url_host", urlHostValidation)
	v.Validator.RegisterValidation("unique", isUniqueField)
}

func (v *Validator) RegisterTranslations() {
	_ = v.Validator.RegisterTranslation("required", v.Translator, func(ut ut.Translator) error {
		return ut.Add("required", "polje je obavezno.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required")
		return t
	})

	_ = v.Validator.RegisterTranslation("email", v.Translator, func(ut ut.Translator) error {
		return ut.Add("email", "polje mora biti email.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email")
		return t
	})

	_ = v.Validator.RegisterTranslation("max", v.Translator, func(ut ut.Translator) error {
		return ut.Add("max", "polje ima limit od maksimalno {0} karaktera.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("max", fe.Param())
		return t
	})

	_ = v.Validator.RegisterTranslation("min", v.Translator, func(ut ut.Translator) error {
		return ut.Add("min", "polje mora imati najmanje {0} karaktera.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Param())
		return t
	})

	_ = v.Validator.RegisterTranslation("receipt_url", v.Translator, func(ut ut.Translator) error {
		return ut.Add("receipt_url", "polje ne predstavlja validan url računa.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("receipt_url")
		return t
	})

	_ = v.Validator.RegisterTranslation("url", v.Translator, func(ut ut.Translator) error {
		return ut.Add("url", "polje nije validan url.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url")
		return t
	})

	_ = v.Validator.RegisterTranslation("url_query_params", v.Translator, func(ut ut.Translator) error {
		return ut.Add("url_query_params", "polju nedostaju obavezni query parametri: {0}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url_query_params", fe.Param())
		return t
	})

	_ = v.Validator.RegisterTranslation("url_host", v.Translator, func(ut ut.Translator) error {
		return ut.Add("url_host", "polje nema validan url host. Validan host je: {0}.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url_host", fe.Param())
		return t
	})

	_ = v.Validator.RegisterTranslation("unique", v.Translator, func(ut ut.Translator) error {
		return ut.Add("unique", "polje mora biti jedinstveno.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique")
		return t
	})
}

func (v *Validator) RegisterTagNameFunc(f validator.TagNameFunc) {
	v.Validator.RegisterTagNameFunc(f)
}

func (v *Validator) RegisterTranslation(tag string, message string, params ...string) {
	v.Validator.RegisterTranslation(
		tag,
		v.Translator,
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

// Custom validation functions

/************************************************************/
/* Custom validation functions								*/
/************************************************************/

// Receipts
func receiptUrlValidation(fl validator.FieldLevel) bool {
	const FISCALIZATION_SYSTEM_HOST = "suf.purs.gov.rs"
	url, err := url.Parse(fl.Field().String())
	if err != nil {
		return false
	}

	hostname := strings.TrimPrefix(url.Hostname(), "www.")

	return hostname == FISCALIZATION_SYSTEM_HOST
}

func urlQueryParamsValidation(fl validator.FieldLevel) bool {
	// Get the URL string from the field
	rawUrl := fl.Field().String()

	// Parse the URL
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return false
	}

	// Get the query parameters from the parsed URL
	queryParams := parsedUrl.Query()
	// Get the required query parameters from the validation tag
	requiredQueryParams := strings.Split(fl.Param(), ",")

	for _, param := range requiredQueryParams {
		if _, exists := queryParams[param]; !exists {
			return false
		}
	}

	return true
}

func urlHostValidation(fl validator.FieldLevel) bool {
	// Get the URL string from the field
	rawUrl := fl.Field().String()

	// Parse the URL
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return false
	}

	// Check for host
	expectedHostname := fl.Param()
	// Actual hostname with trimmed www part
	actualHostname := strings.TrimPrefix(parsedUrl.Hostname(), "www")

	return expectedHostname == actualHostname
}

// Example:
//
//	type RegisterUser struct {
//		Email: unique:users.email,
//		...
//	}
func isUniqueField(fl validator.FieldLevel) bool {
	// struct field value
	value := fl.Field().String()

	// parts[0] - table name
	// parts[1] - table field
	parts := strings.Split(fl.Param(), ".")
	table := parts[0]
	field := parts[1]

	var count int
	if err := db.Instance.Raw("SELECT COUNT(*) FROM `"+table+"` WHERE `"+field+"` = ?", value).Scan(&count).Error; err != nil {
		log.Printf("Error performing unique validation: %s\n", err)
		os.Exit(1)
	}

	return count == 0
}
