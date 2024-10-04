package validation

import (
    "strings"

    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()

    // Register a custom function for the iso639_1 tag
    validate.RegisterValidation("iso639_1", validateISO639_1)
}

func ValidateStruct(s interface{}) error {
    return validate.Struct(s)
}

func validateISO639_1(fl validator.FieldLevel) bool {
    // Later we might want to check against a comprehensive list of ISO 639-1 codes.
    code := fl.Field().String()
    return len(code) == 2 && code == strings.ToLower(code)
}
