package main

import (
	"fmt"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func deburr(source string) (string, error) {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(transformer, source)
	if err != nil {
		fmt.Printf("Error normalizing username: %v\n", err)
		return source, err
	}
	return output, nil
}

func defaults(main string, fallback string) string {
	return ternary(len(main) > 0, main, fallback)
}

func ternary(condition bool, yep string, nope string) string {
	if condition {
		return yep
	}
	return nope
}
