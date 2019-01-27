package add

import (
	"fmt"
	"github.com/dnote/cli/testutils"
	"testing"
)

func TestValidateName(t *testing.T) {
	testCases := []struct {
		input    string
		expected error
	}{
		{
			input:    "trash",
			expected: ErrBookNameReserved,
		},
		{
			input:    "conflicts",
			expected: ErrBookNameReserved,
		},
		{
			input:    "conflict",
			expected: nil,
		},
		{
			input:    "0",
			expected: ErrBookNameInteger,
		},
		{
			input:    "-0",
			expected: ErrBookNameInteger,
		},
		{
			input:    "+0",
			expected: ErrBookNameInteger,
		},
		{
			input:    "11",
			expected: ErrBookNameInteger,
		},
		{
			input:    "-11",
			expected: ErrBookNameInteger,
		},
		{
			input:    "+11",
			expected: ErrBookNameInteger,
		},
		{
			input:    "18446744073709551617",
			expected: ErrBookNameInteger,
		},
		{
			input:    "-18446744073709551617",
			expected: ErrBookNameInteger,
		},
		{
			input:    "2019-01",
			expected: nil,
		},
		{
			input:    "1book",
			expected: nil,
		},
		{
			input:    "book1",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input_%s", tc.input), func(t *testing.T) {
			result := validateName(tc.input)
			testutils.AssertEqual(t, result, tc.expected, "result mismatch")
		})
	}
}
