package helper

import (
	"errors"
	"regexp"
)

func ValidPinfl(pinfl string) error {
	if pinfl == "" {
		return errors.New("error application passport_pinfl requirement body to model")
	}
	pattern := regexp.MustCompile(`^([0-9]{14})$`)

	if !(pattern.MatchString(pinfl)) {
		return errors.New("passport_pinfl must be 14 digits")
	}
	return nil
}

func ValidPassportNumber(number string) error {
	if number == "" {
		return errors.New("error application passport_number requirement body to model")
	}
	pattern := regexp.MustCompile(`^([0-9]{7})$`)

	if !(pattern.MatchString(number)) {
		return errors.New("passport_number must be 7 digits")
	}
	return nil
}

// IsValidPhone ...
func IsValidPhone(phone string) bool {
	r := regexp.MustCompile(`^\+998[0-9]{2}[0-9]{7}$`)
	return r.MatchString(phone)
}

// IsValidEmail ...
func IsValidEmail(email string) bool {
	r := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return r.MatchString(email)
}

// IsValidLogin ...
func IsValidLogin(login string) bool {
	r := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{5,29}$`)
	return r.MatchString(login)
}

// IsValidUUID ...
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func IsValidUUIDV1(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

// IsValidPrice ...
func IsValidPrice(price string) bool {
	r := regexp.MustCompile(`^\d+$`)
	return r.MatchString(price)
}
