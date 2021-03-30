package hw02_unpack_string // nolint:golint,stylecheck

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString    = errors.New("invalid stringg")
	InvalidStringRegexp = `^\d.|[^\\]\d{2,}`
	IsEscaped           bool
)

func Unpack(inStr string) (string, error) {
	if len(inStr) == 0 {
		return inStr, nil
	}

	if ok, _ := regexp.MatchString(InvalidStringRegexp, inStr); ok {
		return inStr, ErrInvalidString
	}

	var (
		outStr     strings.Builder
		next       rune
		inStrRunes = []rune(inStr)
	)

	for i := range inStrRunes[:len(inStrRunes)-1] {
		current := inStrRunes[i]
		next = inStrRunes[i+1]
		if isEscape(current) {
			IsEscaped = true
			continue
		}
		if (unicode.IsDigit(next) && next != 0 && !IsEscaped) || unicode.IsDigit(next) && IsEscaped {
			repeatCount, err := strconv.Atoi(string(next))
			if err != nil {
				return "", err
			}
			outStr.WriteString(strings.Repeat(string(current), repeatCount))
			if IsEscaped {
				IsEscaped = false
			}
			continue
		}
		if !unicode.IsDigit(current) || IsEscaped {
			outStr.WriteRune(current)
			IsEscaped = false
			continue
		}
	}

	if !unicode.IsDigit(next) || IsEscaped {
		outStr.WriteRune(next)
	}

	return outStr.String(), nil
}

func isEscape(val rune) bool {
	return val == '\\' && !IsEscaped
}
