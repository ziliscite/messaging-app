package domain

import (
	"fmt"
)

func ErrEmpty(field string) error {
	return fmt.Errorf("%s cannot be empty", field)
}

func ErrTooLong(field string, maxLength int) error {
	return fmt.Errorf("%s is too long (maximum %d characters)", field, maxLength)
}

func ErrTooShort(field string, minLength int) error {
	return fmt.Errorf("%s is too short (minimum %d characters)", field, minLength)
}

func ErrInvalid(field string) error {
	return fmt.Errorf("%s is invalid", field)
}

func ErrInvalidFormat(field, expectedFormat string) error {
	return fmt.Errorf("%s has invalid format (expected: %s)", field, expectedFormat)
}

func ErrOutOfRange(field string, min, max int) error {
	return fmt.Errorf("%s must be between %d and %d", field, min, max)
}
