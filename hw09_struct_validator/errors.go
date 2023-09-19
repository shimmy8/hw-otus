package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field: %s, error: %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	builder.WriteString("Validation Errors:\n")
	for _, err := range v {
		builder.WriteString(err.Error())
	}
	return builder.String()
}

var (
	ErrInvalidArg        = errors.New("invalid argument: must be struct")
	ErrInvalidRule       = errors.New("invalid rule format, must be key:value")
	ErrInvalidRuleIntVal = errors.New("invalid rule value, must be integer")
	ErrUnknnownRule      = errors.New("unknown rule")
	ErrValidateField     = errors.New("validation error")
	ErrValidateLen       = errors.New("invalid value length")
	ErrValidateRegex     = errors.New("value does not match regexp")
	ErrValidateIn        = errors.New("values is not from set")
	ErrValidateMin       = errors.New("value is smaller than min")
	ErrValidateMax       = errors.New("value is bigger than max")
)

func NewErrUnknnownRule(ruleName string) error {
	return fmt.Errorf("%w: %s", ErrUnknnownRule, ruleName)
}

func NewErrValidateLen(requiredLen int) error {
	return fmt.Errorf("%w, %w: %d required", ErrValidateField, ErrValidateLen, requiredLen)
}

func NewErrValidateRegex(regex string) error {
	return fmt.Errorf("%w, %w: %s", ErrValidateField, ErrValidateRegex, regex)
}

func NewErrValidateIn(in string) error {
	return fmt.Errorf("%w, %w: %s", ErrValidateField, ErrValidateIn, in)
}

func NewErrValidateMin(min int) error {
	return fmt.Errorf("%w, %w: %d", ErrValidateField, ErrValidateMin, min)
}

func NewErrValidateMax(min int) error {
	return fmt.Errorf("%w, %w: %d", ErrValidateField, ErrValidateMax, min)
}
