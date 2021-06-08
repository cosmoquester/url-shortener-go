package utils

import (
	"math/rand"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numLetters = len(letters)

// GenerateRandomURL 은 n 길이의 랜덤 ShortURL을 생성합니다.
func GenerateRandomURL(n uint) string {
	var result string

	for ; n > 0; n-- {
		result = result + string(letters[rand.Intn(numLetters)])
	}
	return result
}

// ValidateURL 은 입력된 문자열이 웹 URL이 맞는지를 체크합니다.
func ValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
