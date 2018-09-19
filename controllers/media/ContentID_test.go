package media

import (
	"math/rand"
	"testing"
)

func TestBuildContentID(t *testing.T) {

	maxLength := randString(maxContentIDLength)
	greaterThanMaxLength := randString(maxContentIDLength + 1)
	significantlyGreaterThanMaxLength := randString(maxContentIDLength + 51231)

	var tests = []struct {
		input       string
		expected    contentID
		errorResult error
	}{
		{"", "", nil},
		{"a", "a", nil},
		{"MIMETYPE", "MIMETYPE", nil},
		{"1", "1", nil},
		{"12310239123", "12310239123", nil},
		{maxLength, contentID(maxLength), nil},
		{greaterThanMaxLength, "", &ContentIDLengthError{len(greaterThanMaxLength)}},
		{significantlyGreaterThanMaxLength, "", &ContentIDLengthError{len(significantlyGreaterThanMaxLength)}},
	}

	for num, test := range tests {
		actual, err := NewContentID(test.input)
		if err != test.errorResult {
			if err != nil && test.errorResult != nil && err.Error() == test.errorResult.Error() {
				continue
			}
			t.Errorf("Error for test %v returned different from expected. (%v) expected, (%v) received", num, test.errorResult, err)
		}
		if test.expected != actual {
			t.Errorf("Result for test %v returned different from expected. (%v) expected, (%v) received", num, test.expected, actual)
		}
	}
	t.Log("all tests completed")
}

func randString(length int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	result := make([]rune, length)

	for index := 0; index < length; index++ {
		result[index] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}
