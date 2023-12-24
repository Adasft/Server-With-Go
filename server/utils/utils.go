package utils

import (
	"net/mail"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsEmptyStr(str string) bool {
	return len(str) == 0
}

func Append[T comparable](slice *[]T, values ...T) {
	*slice = append(*slice, values...)
}
