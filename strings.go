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
