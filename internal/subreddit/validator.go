package subreddit

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	SubredditNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

const (
	ErrSubredditNameRequired = "subreddit name is required"
	ErrSubredditNameTooShort = "subreddit name must be at least %d characters"
	ErrSubredditNameTooLong  = "subreddit name must be at most %d characters"
	ErrSubredditNameInvalid  = "subreddit name can only contain letters, numbers, and underscores"
	ErrSubredditNameTaken    = "subreddit name already taken"

	ErrDisplayNameRequired      = "display name is required"
	ErrDisplayNameNoWhitespaces = "display name cannot have leading or trailing whitespace"
	ErrDisplayNameTooLong       = "display name must be at most %d characters"

	ErrDescriptionTooLong = "description must be at most %d characters"
	ErrIconURLTooLong     = "icon URL must be at most %d characters"

	NameMinLen        = 3
	NameMaxLen        = 21
	DisplayNameMaxLen = 255
	DescriptionMaxLen = 500
	IconURLMaxLen     = 500
)

type Validator struct {
	repo      *Repository
	nameRegex *regexp.Regexp
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo:      repo,
		nameRegex: SubredditNameRegex,
	}
}

func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

func (ve ValidationErrors) Error() string {
	return "validation failed"
}

func (v *Validator) ValidateNameFormat(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New(ErrSubredditNameRequired)
	}

	if len(name) < NameMinLen {
		return errors.New(fmt.Sprintf(ErrSubredditNameTooShort, NameMinLen))
	}

	if len(name) > NameMaxLen {
		return errors.New(fmt.Sprintf(ErrSubredditNameTooLong, NameMaxLen))
	}

	if !v.nameRegex.MatchString(name) {
		return errors.New(ErrSubredditNameInvalid)
	}

	return nil
}

func (v *Validator) ValidateDisplayNameFormat(displayName string) error {
	displayName = strings.TrimSpace(displayName)

	if displayName == "" {
		return errors.New(ErrDisplayNameRequired)
	}

	if displayName != strings.TrimSpace(displayName) {
		return errors.New(ErrDisplayNameNoWhitespaces)
	}

	if len(displayName) > DisplayNameMaxLen {
		return errors.New(fmt.Sprintf(ErrDisplayNameTooLong, DisplayNameMaxLen))
	}

	return nil
}

func (v *Validator) ValidateDescriptionFormat(description *string) error {
	if description == nil {
		return nil // Optional field
	}

	desc := strings.TrimSpace(*description)
	if len(desc) > DescriptionMaxLen {
		return errors.New(fmt.Sprintf(ErrDescriptionTooLong, DescriptionMaxLen))
	}

	return nil
}

func (v *Validator) ValidateIconURLFormat(iconURL *string) error {
	if iconURL == nil {
		return nil // Optional field
	}

	url := strings.TrimSpace(*iconURL)
	if len(url) > IconURLMaxLen {
		return errors.New(fmt.Sprintf(ErrIconURLTooLong, IconURLMaxLen))
	}

	return nil
}

func (v *Validator) ValidateNameExists(ctx context.Context, name string) error {
	exists, _ := v.repo.ExistsByName(ctx, strings.ToLower(name))
	if exists {
		return errors.New(ErrSubredditNameTaken)
	}
	return nil
}

func (v *Validator) ValidateCreateSubredditInput(
	ctx context.Context,
	name string,
	displayName string,
	description *string,
	iconURL *string,
) ValidationErrors {
	var errs ValidationErrors

	if err := v.ValidateNameFormat(name); err != nil {
		errs = append(errs, NewValidationError("name", err.Error()))
	}

	if err := v.ValidateDisplayNameFormat(displayName); err != nil {
		errs = append(errs, NewValidationError("display_name", err.Error()))
	}

	if err := v.ValidateDescriptionFormat(description); err != nil {
		errs = append(errs, NewValidationError("description", err.Error()))
	}

	if err := v.ValidateIconURLFormat(iconURL); err != nil {
		errs = append(errs, NewValidationError("icon_url", err.Error()))
	}

	if len(errs) == 0 {
		if err := v.ValidateNameExists(ctx, name); err != nil {
			errs = append(errs, NewValidationError("name", err.Error()))
		}
	}

	return errs
}

func (v *Validator) ValidateUpdateSubredditInput(req UpdateSubredditRequest) ValidationErrors {
	var errs ValidationErrors

	if req.DisplayName != nil {
		if err := v.ValidateDisplayNameFormat(*req.DisplayName); err != nil {
			errs = append(errs, NewValidationError("display_name", err.Error()))
		}
	}

	if req.Description != nil {
		if err := v.ValidateDescriptionFormat(req.Description); err != nil {
			errs = append(errs, NewValidationError("description", err.Error()))
		}
	}

	if req.IconURL != nil {
		if err := v.ValidateIconURLFormat(req.IconURL); err != nil {
			errs = append(errs, NewValidationError("icon_url", err.Error()))
		}
	}

	return errs
}
