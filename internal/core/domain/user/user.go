package user

import (
	"fmt"
	"github.com/ziliscite/messaging-app/internal/core/domain"
	"regexp"
	"time"
	"unicode"
)

func ErrPassword(rule string) error {
	return fmt.Errorf("password must %s", rule)
}

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ValidateUsername(username string) (string, error) {
	if username == "" {
		return "", domain.ErrEmpty("username")
	}

	if len(username) < 3 {
		return "", domain.ErrTooShort("username", 3)
	}

	if len(username) > 100 {
		return "", domain.ErrTooLong("username", 100)
	}

	return username, nil
}

func ValidateEmail(email string) (string, error) {
	if email == "" {
		return "", domain.ErrEmpty("email")
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if re.MatchString(email) == false {
		return "", domain.ErrInvalid("email")
	}

	return email, nil
}

func ValidatePassword(password string) (string, error) {
	if len(password) < 8 {
		return "", ErrPassword("at least 8 characters long")
	}

	if len(password) > 72 {
		return "", ErrPassword("not exceed 72 characters")
	}

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

	if !hasUpper {
		return "", ErrPassword("contain at least one uppercase letter")
	}
	if !hasLower {
		return "", ErrPassword("contain at least one lowercase letter")
	}
	if !hasNumber {
		return "", ErrPassword("contain at least one number")
	}
	if !hasSpecial {
		return "", ErrPassword("contain at least one special character")
	}

	return password, nil
}
