package main

import "testing"

func TestGenerateRandomUrl(t *testing.T) {
	var str1, str2 string

	str1 = generateRandomUrl(10)
	str2 = generateRandomUrl(10)

	if str1 == str2 {
		t.Error("random logit is not working!")
	}
}

func TestValidateUrl(t *testing.T) {
	if validateUrl("naver.com") {
		t.Error("url validation failed!")
	} else if !validateUrl("https://a.a") {
		t.Error("url validation failed!")
	} else if validateUrl("http:/hello.com") {
		t.Error("url validation failed!")
	}
}
