package path

import (
	"path/filepath"
	"regexp"
)

func GetPathRegex(src string) (*regexp.Regexp, error) {
	pattern := filepath.Base(src)

	reg, err := regexp.Compile("\\*")
	if err != nil {
		return nil, err
	}

	pattern = "^" + string(reg.ReplaceAll([]byte(pattern), []byte("[\\w|\\W]*"))) + "$"

	return regexp.Compile(pattern)
}

func isFilePathIsRegex(reg string) bool {
	specialChars := "*+?|[]{}()"

	for _, specChar := range specialChars {
		for _, char := range reg {
			if char == specChar {
				return true
			}
		}
	}

	return false
}
