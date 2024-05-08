package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomString(t *testing.T) {
	strLen := 10
	result := GetRandomString(strLen)
	assert.Equal(t, strLen, len(result), "The length of the returned string does not match the expected length")

	// there is a small probability that this test fails even if the function works correctly
	result1 := GetRandomString(strLen)
	result2 := GetRandomString(strLen)
	assert.NotEqual(t, result1, result2, "Generated strings shouldn't be equal")

	for _, runeValue := range result {
		assert.True(t, strings.ContainsRune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", runeValue))
	}
}

func TestGenerateSaltedHash(t *testing.T) {
	salt := "12345"
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", "5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5"},
		{"Simple string", "Hello, World!", "c20a81dd8942ae16134d4c71ad5de156996c73220a3b52ecd651cfff12b21aba"},
		{"Numeric string", "12345", "e4a0a90e5ac07d5435c6f25c4cf7cc565becb797bb5b83c515bc427ef32a4770"},
		{"Special chars string", "!@#$%^&*(){}", "514fe42a1a19254e7bbd5983e296c36fd92a8b8c0bd60fee8f83488e8fc37cdc"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := GenerateSaltedHash(tc.input, salt)
			if got != tc.expected {
				t.Errorf("GenerateHashFromString(%v): expected %v, but got %v",
					tc.input, tc.expected, got)
			}

			assert.Equal(t, 64, len(got), "The length of the generated hash is not 64")
			assert.Equal(t, tc.expected, got, "Hash does not match expected result")
			assert.NotEqual(t, tc.input, got, "The generated hash should not be equal to the original string")
		})
	}
}

func TestGenerateHashFromString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"Simple string", "Hello, World!", "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"},
		{"Numeric string", "12345", "5994471abb01112afcc18159f6cc74b4f511b99806da59b3caf5a9c173cacfc5"},
		{"Special chars string", "!@#$%^&*(){}", "5e5451bc577c6b92a56617311be68d6c8bbe06ec4fe85c404d610a329c462072"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := GenerateHashFromString(tc.input)
			if got != tc.expected {
				t.Errorf("GenerateHashFromString(%v): expected %v, but got %v",
					tc.input, tc.expected, got)
			}

			assert.Equal(t, 64, len(got), "The length of the generated hash is not 64")
			assert.Equal(t, tc.expected, got, "Hash does not match expected result")
			assert.NotEqual(t, tc.input, got, "The generated hash should not be equal to the original string")
		})
	}
}
