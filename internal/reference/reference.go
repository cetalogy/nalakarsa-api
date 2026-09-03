package reference

import (
	"fmt"
	"regexp"
	"strings"
)

var namePattern = regexp.MustCompile(`^[\p{L}\p{N}][\p{L}\p{N} .,&'()/;:+-]{1,254}$`)

func NormalizeName(value string) (string, error) {
	name := strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	if len([]rune(name)) < 2 || len([]rune(name)) > 255 || !namePattern.MatchString(name) {
		return "", fmt.Errorf("name contains invalid characters or length")
	}
	return name, nil
}
