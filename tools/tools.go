package tools

import (
	"fmt"
	"regexp"
)

func IsValidGuess(g string) (bool, error) {
	r, _ := regexp.Compile(`^[A-Za-z]+$`)

	if len(g) != 5 || !r.MatchString(g) {
		return false, fmt.Errorf("incorrect format")
	}

	return true, nil
}

func FormatResponse(result string) string {
	if result == "correct" {
		return "🟩"
	}

	if result == "present" {
		return "🟧"
	}

	if result == "absent" {
		return "⬛"
	}

	return result
}
