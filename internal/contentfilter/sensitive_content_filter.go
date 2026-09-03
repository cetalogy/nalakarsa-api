package contentfilter

import (
	"errors"
	"os"
	"strings"
	"unicode"

	"nalakarsa/internal/common/constant"
)

var ErrSensitiveContent = errors.New("input contains prohibited words")
func Validate(value string) error {
	text := normalize(value)
	if text == "" {
		return nil
	}

	configuredWords := strings.Join([]string{
		constant.DefaultSensitiveWords,
		os.Getenv("SENSITIVE_WORDS"),
	}, ",")
	for _, rawWord := range strings.Split(configuredWords, ",") {
		word := normalize(rawWord)
		if word != "" && containsTerm(text, word) {
			return ErrSensitiveContent
		}
	}
	return nil
}

func containsTerm(text, term string) bool {
	if strings.Contains(term, " ") {
		return strings.Contains(text, term)
	}
	for _, token := range strings.Fields(text) {
		if token == term {
			return true
		}
	}
	return false
}

func normalize(value string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}
