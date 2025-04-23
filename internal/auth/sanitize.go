package auth

import "regexp"

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9_.]+@[a-zA-Z0-9]+\.[a-zA-Z]+$`)
var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{8,20}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func IsValidPassword(password string) bool {
	return passwordRegex.MatchString(password)
}
