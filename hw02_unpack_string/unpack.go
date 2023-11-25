package hw02unpackstring

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	//пустое значение возвращает пустое значение
	if str == "" {
		return "", nil
	}

	runes := []rune(str)
	//не должно начинаться с цифры
	if _, err := strconv.Atoi(string(runes[0])); err == nil {
		return "", ErrInvalidString
	}
	//недопустимы числа
	reg := regexp.MustCompile(`[0-9]{2,}`)
	if reg.MatchString(str) {
		return "", ErrInvalidString
	}

	result := strings.Builder{}
	var r rune
	for i := range runes {
		rn := runes[i]
		if val, err := strconv.Atoi(string(rn)); err == nil {
			result.WriteString(strings.Repeat(string(r), val))
			r = 0
		} else {
			if r != 0 {
				result.WriteRune(r)
			}
			r = rn
		}
	}
	if r != 0 {
		result.WriteRune(r)
	}
	return result.String(), nil
}
