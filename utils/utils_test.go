package utils

import (
	"testing"
)

func TestGenerateRandomURL(t *testing.T) {
	var str1, str2 string

	str1 = GenerateRandomURL(10)
	str2 = GenerateRandomURL(10)

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
		if ValidateURL(testCase.URL) != testCase.IsValid {
			t.Errorf("url validation failed! \"URL: %v\", \"isValid: %v\"", testCase.URL, testCase.IsValid)
		}
	}
}
