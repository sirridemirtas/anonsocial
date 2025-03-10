package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/sirridemirtas/anonsocial/data"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("university", validateUniversity)
}

func validateUniversity(fl validator.FieldLevel) bool {
	return data.IsValidUniversityID(fl.Field().String())
}

func ValidateUser(user interface{}) []string {
	var errors []string
	err := validate.Struct(user)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				errors = append(errors, field+" is required")
			case "alphanum":
				errors = append(errors, field+" must contain only letters and numbers")
			case "min":
				errors = append(errors, field+" must be at least "+err.Param()+" characters long")
			case "max":
				errors = append(errors, field+" must not exceed "+err.Param()+" characters")
			case "university":
				errors = append(errors, "invalid university ID")
			}
		}
	}

	return errors
}

// ValidateUsername checks if a username meets the requirements:
// - alphanumeric only
// - between 3-16 characters long
// Returns a slice of error messages, empty if valid
func ValidateUsername(username string) []string {
	var errors []string

	// Check if empty
	if username == "" {
		errors = append(errors, "username is required")
		return errors
	}

	// Check length (3-16 characters)
	length := utf8.RuneCountInString(username)
	if length < 3 {
		errors = append(errors, "username must be at least 3 characters long")
	}
	if length > 16 {
		errors = append(errors, "username must not exceed 16 characters")
	}

	// Check if alphanumeric only
	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumeric.MatchString(username) {
		errors = append(errors, "username must contain only letters and numbers")
	}

	return errors
}
