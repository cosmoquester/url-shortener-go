package main

import (
	"math/rand"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numLetters = len(letters)

func generateRandomUrl(n uint) string {
	var result string

	for ; n > 0; n-- {
		result = result + string(letters[rand.Intn(numLetters)])
	}
	return result
}

func validateUrl(url string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}
