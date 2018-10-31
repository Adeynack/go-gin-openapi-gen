package gen

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitByOneOf(t *testing.T) {
	for _, c := range []struct {
		input      string
		separators []string
		expected   []string
	}{
		{
			input:      "abcdefghijklmnopqrstuvwxyz",
			separators: []string{"a", "e", "i", "o", "u", "y"},
			expected:   []string{"", "bcd", "fgh", "jklmn", "pqrst", "vwx", "z"},
		},
		{
			input:      "this.is_supposed_to-split-this.sentence",
			separators: []string{".", "-", "_"},
			expected:   []string{"this", "is", "supposed", "to", "split", "this", "sentence"},
		},
	} {
		t.Run(fmt.Sprintf("%q split by %v", c.input, c.separators), func(t *testing.T) {
			assert.Equal(t, c.expected, splitByOneOf(c.input, c.separators...))
		})
	}
}

func TestToGoFieldName(t *testing.T) {
	for _, c := range []struct {
		input    string
		expected string
	}{
		{"id", "ID"},
		{"name", "Name"},
		{"owner_id", "OwnerID"},
	} {
		t.Run(c.input, func(t *testing.T) {
			actual := toGoFieldName(c.input)
			assert.Equal(t, c.expected, actual)
		})
	}
}
