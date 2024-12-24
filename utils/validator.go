package utils

import (
	"strings"

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
