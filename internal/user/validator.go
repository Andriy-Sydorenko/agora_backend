package user

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// Separating simple validators are a little bit of overhead in my case, but feels more readable
var (
	EmailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	PasswordRegex = regexp.MustCompile(`[A-Z].*[a-z].*[0-9]|[A-Z].*[0-9].*[a-z]|[a-z].*[A-Z].*[0-9]|[a-z].*[0-9].*[A-Z]|[0-9].*[A-Z].*[a-z]|[0-9].*[a-z].*[A-Z]`) // FIXME: rewrite this silly regexp
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	//v.RegisterValidation("not_common", validateNotCommon)
	return &Validator{validate: v}
}

func (v *Validator) ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return errors.New("username must be 3-50 characters")
	}
	if !UsernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func (v *Validator) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !PasswordRegex.MatchString(password) {
		return errors.New("password must contain uppercase, lowercase, and number")
	}

	return nil
}

func (v *Validator) ValidateEmail(email string) error {
	if !EmailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}
