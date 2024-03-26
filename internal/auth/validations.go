package auth

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/larry-lw-chan/goti/internal/utils/translate"
)

/************************************************************
* Validations
***********************************************************/

// Create User Validation
type CreateUserValidation struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func validateCreateUser(r *http.Request) (errs []error) {
	// Check if passwords match
	if r.FormValue("password") != r.FormValue("confirm_password") {
		errs = append(errs, errors.New("passwords do not match"))
	}

	// Check if privacy policy is agreed
	if r.FormValue("privacy") == "off" {
		errs = append(errs, errors.New("please agree to the privacy policy"))
	}

	// Validate User Input
	validate := validator.New()
	createUser := CreateUserValidation{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	vErrs := validate.Struct(&createUser)
	if vErrs != nil {
		errs = append(errs, translate.Errors(vErrs, validate)...)
	}

	// Returns nil if no errors
	return errs
}

// Login Validation
type LoginUserValidation struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func validateLoginUser(r *http.Request) (errs []error) {
	// Validate User Input
	validate := validator.New()

	loginUser := LoginUserValidation{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	vErrs := validate.Struct(&loginUser)
	if vErrs != nil {
		errs = append(errs, translate.Errors(vErrs, validate)...)
	}

	// Returns nil if no errors
	return errs
}
