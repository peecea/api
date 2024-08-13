package utils

import "regexp"

/*
 CHECK IF MOBILE PHONE NUMBER IS NUMERIC AND A LENGHT OF THIRTEEN
*/
func IsValidPhone(phone string) bool {
	pattern := `^[+]+\d{12}$`

	re := regexp.MustCompile(pattern)

	return re.MatchString(phone)
}