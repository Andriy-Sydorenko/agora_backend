package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/Andriy-Sydorenko/agora_backend/internal/user"
	"regexp"
	"strings"
)

var (
	EmailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	PasswordRegex = regexp.MustCompile(`[A-Z].*[a-z].*[0-9]|[A-Z].*[0-9].*[a-z]|[a-z].*[A-Z].*[0-9]|[a-z].*[0-9].*[A-Z]|[0-9].*[A-Z].*[a-z]|[0-9].*[a-z].*[A-Z]`) // FIXME: rewrite this silly regexp
)

const (
	ErrEmailRequired         = "email is required"
	ErrEmailInvalid          = "invalid email format"
	ErrEmailNoWhitespaces    = "email cannot have leading or trailing whitespace"
	ErrEmailTaken            = "email already registered"
	ErrUsernameRequired      = "username is required"
	ErrUsernameNoWhitespaces = "username cannot have leading or trailing whitespace"
	ErrUsernameTooShort      = "username must be at least %d characters"
	ErrUsernameTooLong       = "username must be at most %d characters"
	ErrUsernameInvalid       = "username can only contain letters, numbers, and underscores"
	ErrUsernameTaken         = "username already taken"
	ErrPasswordRequired      = "password is required"
	ErrPasswordNoWhitespaces = "password cannot have leading or trailing whitespace"
	ErrPasswordTooShort      = "password must be at least %d characters"
	ErrPasswordTooLong
	ErrPasswordWeak = "password must contain uppercase, lowercase, and number"

	UsernameMinLen = 3
	UsernameMaxLen = 50
	PasswordMinLen = 8
	PasswordMaxLen = 30
)

type Validator struct {
	userService   *user.Service
	emailRegex    *regexp.Regexp
	usernameRegex *regexp.Regexp
	passwordRegex *regexp.Regexp
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func NewValidator(userService *user.Service) *Validator {
	return &Validator{
		userService:   userService,
		emailRegex:    EmailRegex,
		usernameRegex: UsernameRegex,
		passwordRegex: PasswordRegex,
	}
}

func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// Implementing Error method to convert ValidationErrors to error type interface and follow Go's error handling contract
// FIXME: come up with another logic, this is straight up a workaround
func (ve ValidationErrors) Error() string {
	return "validation failed"
}

func (v *Validator) ValidateEmailFormat(email string) error {
	if email == "" {
		return errors.New(ErrEmailRequired)
	}
	if email != strings.TrimSpace(email) {
		return errors.New(ErrEmailNoWhitespaces)
	}
	if !v.emailRegex.MatchString(email) {
		return errors.New(ErrEmailInvalid)
	}
	return nil
}

func (v *Validator) ValidateUsernameFormat(username string) error {
	if username == "" {
		return errors.New(ErrUsernameRequired)
	}
	if username != strings.TrimSpace(username) {
		return errors.New(ErrUsernameNoWhitespaces)
	}
	if len(username) < UsernameMinLen {
		return errors.New(fmt.Sprintf(ErrUsernameTooShort, UsernameMinLen))
	}
	if len(username) > UsernameMaxLen {
		return errors.New(fmt.Sprintf(ErrUsernameTooLong, UsernameMaxLen))
	}
	if !v.usernameRegex.MatchString(username) {
		return errors.New(ErrUsernameInvalid)
	}
	return nil
}

func (v *Validator) ValidatePasswordFormat(password string) error {
	if password == "" {
		return errors.New(ErrPasswordRequired)
	}
	if password != strings.TrimSpace(password) {
		return errors.New(ErrPasswordNoWhitespaces)
	}
	if len(password) < PasswordMinLen {
		return errors.New(fmt.Sprintf(ErrPasswordTooShort, PasswordMinLen))
	}
	if len(password) > PasswordMaxLen {
		return errors.New(fmt.Sprintf(ErrPasswordTooLong, PasswordMaxLen))
	}

	if !v.passwordRegex.MatchString(password) {
		return errors.New(ErrPasswordWeak)
	}
	return nil
}

func (v *Validator) ValidateEmailAvailability(ctx context.Context, email string) error {
	alreadyExists, _ := v.userService.ExistsByEmail(ctx, email)
	if alreadyExists {
		return errors.New(ErrEmailTaken)
	}
	return nil
}

func (v *Validator) ValidateUsernameAvailability(ctx context.Context, username string) error {
	alreadyExists, _ := v.userService.ExistsByUsername(ctx, username)
	if alreadyExists {
		return errors.New(ErrUsernameTaken)
	}
	return nil
}

func (v *Validator) ValidateRegistrationInput(ctx context.Context, email, username, password string) ValidationErrors {
	var errs ValidationErrors
	// Format validation
	if err := v.ValidateEmailFormat(email); err != nil {
		errs = append(errs, NewValidationError("email", err.Error()))
	}
	if err := v.ValidateUsernameFormat(username); err != nil {
		errs = append(errs, NewValidationError("username", err.Error()))
	}
	if err := v.ValidatePasswordFormat(password); err != nil {
		errs = append(errs, NewValidationError("password", err.Error()))
	}

	// Business validation(expensive by DB hits, performed only if formatting validation succeeds)
	if len(errs) == 0 {
		if err := v.ValidateEmailAvailability(ctx, email); err != nil {
			errs = append(errs, NewValidationError("email", err.Error()))
		}
		if err := v.ValidateUsernameAvailability(ctx, username); err != nil {
			errs = append(errs, NewValidationError("username", err.Error()))
		}
	}

	return errs
}
