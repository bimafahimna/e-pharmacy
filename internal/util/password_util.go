package util

import (
	"regexp"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
)

var (
	hasLower   = regexp.MustCompile(`[a-z]`).MatchString
	hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
	hasDigit   = regexp.MustCompile(`[0-9]`).MatchString
	hasSpecial = regexp.MustCompile(`[\W]`).MatchString
)

func ValidatePassword(password string) error {
	if !hasLower(password) || !hasUpper(password) || !hasDigit(password) || !hasSpecial(password) {
		return apperror.ErrInvalidPasswordFormat
	}

	return nil
}
