package general_functions

import (
	"regexp"
	"strings"
)

// RemoveNonAlphanumeric https://golangcode.com/how-to-remove-all-non-alphanumerical-characters-from-a-string/
func RemoveNonAlphanumeric(str string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}

	return strings.ToLower(reg.ReplaceAllString(str, "-")), nil
}
