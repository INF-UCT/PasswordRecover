package handler

import (
	"strings"
	"unicode"
)

const passwordMinLength = 8

const passwordSpecialSet = "!@#$%^&*()_+-=[]{};':\"\\|,.<>/?~`"

type fieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func validatePassword(pw string) []fieldError {
	var errs []fieldError

	if utf8Len(pw) < passwordMinLength {
		errs = append(errs, fieldError{
			Field:   "new_password",
			Code:    "too_short",
			Message: "Debe tener al menos 8 caracteres.",
		})
	}

	hasLower, hasUpper, hasDigit, hasSpecial := false, false, false, false
	for _, r := range pw {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case strings.ContainsRune(passwordSpecialSet, r):
			hasSpecial = true
		}
	}

	if !hasLower {
		errs = append(errs, fieldError{
			Field:   "new_password",
			Code:    "missing_lowercase",
			Message: "Debe incluir al menos una letra minúscula.",
		})
	}
	if !hasUpper {
		errs = append(errs, fieldError{
			Field:   "new_password",
			Code:    "missing_uppercase",
			Message: "Debe incluir al menos una letra mayúscula.",
		})
	}
	if !hasDigit {
		errs = append(errs, fieldError{
			Field:   "new_password",
			Code:    "missing_digit",
			Message: "Debe incluir al menos un número.",
		})
	}
	if !hasSpecial {
		errs = append(errs, fieldError{
			Field:   "new_password",
			Code:    "missing_symbol",
			Message: "Debe incluir al menos un símbolo especial.",
		})
	}

	return errs
}

func validatePasswordMatch(a, b string) *fieldError {
	if a != b {
		return &fieldError{
			Field:   "confirmed_password",
			Code:    "password_mismatch",
			Message: "Las contraseñas no coinciden.",
		}
	}
	return nil
}

func utf8Len(s string) int {
	return len([]rune(s))
}
