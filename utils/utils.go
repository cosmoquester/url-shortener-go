package utils

import (
	"math/rand"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numLetters = len(letters)

func GenerateRandomURL(n uint) string {
	var result string

	for ; n > 0; n-- {
		result = result + string(letters[rand.Intn(numLetters)])
	}
	return result
}

func ValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
