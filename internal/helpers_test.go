package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"\tquery\t", "query"},
		{"\nnewline\n", "newline"},
	}

	for _, test := range tests {
		result := sanitizeQuery(test.input)
		assert.Equal(t, test.expected, result)
	}
}
