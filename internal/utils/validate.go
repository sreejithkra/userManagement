package utils

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/go-playground/validator"
)

func validateNameOrInitials(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	initialsRegex := regexp.MustCompile(`^([A-Z]\. )*[A-Z]\.$`)
	nameRegex := regexp.MustCompile(`^[A-Za-z]+$`)

	return initialsRegex.MatchString(value) || nameRegex.MatchString(value)
}
// func passwordValidation(fl validator.FieldLevel) bool {
// 	password := fl.Field().String()

// 	if match, _ := regexp.MatchString(`[A-Z]`, password); !match {
// 		return false
// 	}

// 	if match, _ := regexp.MatchString(`[a-z]`, password); !match {
// 		return false
// 	}

// 	if match, _ := regexp.MatchString(`[0-9]`, password); !match {
// 		return false
// 	}

// 	if match, _ := regexp.MatchString(`[!@#\$%\^&\*]`, password); !match {
// 		return false
// 	}

// 	return true
// }
func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
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

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return false
	}
	return true
}

func Validate(data interface{}) error {
	validate := validator.New()

	validate.RegisterValidation("nameOrInitials", validateNameOrInitials)
	validate.RegisterValidation("password", passwordValidation)

	err := validate.Struct(data)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {

			switch e.Tag() {

			case "required":
				return fmt.Errorf("%s is required", e.Field())

			case "email":
				return fmt.Errorf("%s is not a valid email address", e.Field())

			case "numeric":
				return fmt.Errorf("%s shouls contain only digits", e.Field())

			case "len":
				return fmt.Errorf("%s shouls have a length of %s", e.Field(), e.Param())

			case "min":
				return fmt.Errorf("%s shouls have a minimum length of %s", e.Field(), e.Param())

			case "nameOrInitials":
				return fmt.Errorf("%s should be either initials or a regular name", e.Field())

			case "password":
				return fmt.Errorf("%s should contain at least one uppercase letter, one lowercase letter, one digit, and one special character", e.Field())

			case "max":
				return fmt.Errorf("%s exceeds the maximum length", e.Field())

			case "alpha":
				return fmt.Errorf("%s should contain only alphabetic characters", e.Field())

			case "gt":
				return fmt.Errorf("%s must be greater than zero", e.Field())

			default:
				return fmt.Errorf("validation error for field %s", e.Field())

			}
		}
	}
	return nil
}
