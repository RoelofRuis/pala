package pala

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// ParseInt is a literal evaluator for integers in string representation.
func ParseInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

// ParseString is a literal evaluator for plain strings.
func ParseString(s string) (string, error) {
	return s, nil
}

// ParseQuotedString is a literal evaluator for strings using double quotes.
func ParseQuotedString(s string) (string, error) {
	match, _ := regexp.MatchString("^\"[^\"]*\"$", s)
	if !match {
		return "", fmt.Errorf("no valid quoted string")
	}
	return strings.Trim(s, "\""), nil
}

// ParseRational is a literal evaluator for rationals.
func ParseRational(s string) (*big.Rat, error) {
	match, _ := regexp.MatchString("^[0-9]+/[0-9]+$", s)
	if !match {
		return nil, fmt.Errorf("no valid rational")
	}

	parts := strings.Split(s, "/")
	numerator, _ := strconv.ParseInt(parts[0], 10, 64)
	denominator, _ := strconv.ParseInt(parts[1], 10, 64)
	if denominator == 0 {
		numerator = 0
		denominator = 1
	}

	return big.NewRat(numerator, denominator), nil
}
