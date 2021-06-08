package main

import (
	"testing"
)

func TestGenerateRandomURL(t *testing.T) {
	var str1, str2 string

	str1 = generateRandomURL(10)
	str2 = generateRandomURL(10)

	if str1 == str2 {
		t.Error("random logit is not working!")
	}
}

func TestValidateURL(t *testing.T) {
	cases := []struct {
		URL     string
		IsValid bool
	}{
		{"naver.com", false},
		{"https://a.a", true},
		{"http:/hello.com", false},
	}

	for _, testCase := range cases {
		if validateURL(testCase.URL) != testCase.IsValid {
			t.Error("url validation failed!")
		}
	}
}
