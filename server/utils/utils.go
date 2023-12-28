package utils

import (
	"net/mail"
	"regexp"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPhoneNumber(phone string) bool {
	pattern := `^\d{10}$`
	regExp := regexp.MustCompile(pattern)

	return regExp.MatchString(phone)
}

func IsEmptyStr(str string) bool {
	return len(str) == 0
}

func Append[T comparable](slice *[]T, values ...T) {
	*slice = append(*slice, values...)
}
