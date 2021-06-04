package main

import "testing"

func TestGenerateRandomURL(t *testing.T) {
	var str1, str2 string

	str1 = generateRandomURL(10)
	str2 = generateRandomURL(10)

	if str1 == str2 {
		t.Error("random logit is not working!")
	}
}

func TestValidateURL(t *testing.T) {
	if validateURL("naver.com") {
		t.Error("url validation failed!")
	} else if !validateURL("https://a.a") {
		t.Error("url validation failed!")
	} else if validateURL("http:/hello.com") {
		t.Error("url validation failed!")
	}
}
