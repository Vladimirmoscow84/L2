package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUnpacking(t *testing.T) {
	tests := []struct {
		name     string
		string1  string
		expected string
	}{
		{
			name:     "1",
			string1:  "a5\\43bfg2\\5",
			expected: "aaaaa444bfgg5",
		},
		{
			name:     "2",
			string1:  "abcdef",
			expected: "abcdef",
		},
		{
			name:     "3",
			string1:  "\\13\\42\\6\\5\\7",
			expected: "11144657",
		},
		{
			name:     "4",
			string1:  "23",
			expected: "",
		},
		{
			name:     "5",
			string1:  "qwe\\45",
			expected: "qwe44444",
		},
		{
			name:     "6",
			string1:  "45",
			expected: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, _ := getUnpacking(test.string1)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func BenchmarkGetUnpacking(b *testing.B) {
	string1 := "s4fghytr6\\5\\67hg5"
	for b.Loop() {
		getUnpacking(string1)
	}
}
