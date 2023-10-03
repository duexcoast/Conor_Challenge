package challenge

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestCCValidate(t *testing.T) {
	tt := map[string]struct {
		testID   int
		input    string
		expected bool
	}{
		"Four consecutive repeated digits": {
			testID:   0,
			input:    "4444-5555-6666-7777",
			expected: false,
		},
		"Four consecutive digits across hyphen boundary": {
			testID:   1,
			input:    "5044-4420-8173-7744",
			expected: false,
		},
		"Valid input": {
			testID:   2,
			input:    "6345093243215543",
			expected: true,
		},
		"Valid input with hyphens": {
			testID:   3,
			input:    "4908-1573-6339-2872",
			expected: true,
		},
		"Too short": {
			testID:   4,
			input:    "511898531667435",
			expected: false,
		},
		"Too long": {
			testID:   5,
			input:    "6432-9099-4325-58932",
			expected: false,
		},

		"Out of place hyphen": {
			testID:   6,
			input:    "54910-987-2578-3904",
			expected: false,
		},
		"Trailing hyphen": {
			testID:   7,
			input:    "5491-0987-2578-3904-",
			expected: false,
		},
		"Leading hyphen": {
			testID:   8,
			input:    "-5491-0987-2578-3904",
			expected: false,
		},
		"Contains non-digit characters": {
			testID:   9,
			input:    "6O438I4033277914",
			expected: false,
		},
		"Starting character that is not [456]": {
			testID:   10,
			input:    "3081-2855-9842-8003",
			expected: false,
		},
	}

	re := regexp.MustCompile(ccRegExp)

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Logf("\tTest %d:\t%s", test.testID, name)

			got := validateCC(re, test.input)
			if got != test.expected {
				t.Logf("\t\tTest %d:\tExp: %v", test.testID, test.expected)
				t.Logf("\t\tTest %d:\tGot: %v", test.testID, got)
				if test.expected {
					t.Fatalf("\t%s\tTest %d: Should have correctly validated the input string", failed, test.testID)
				} else {
					t.Fatalf("\t%s\tTest %d: Should have correctly *invalidated* the input string", failed, test.testID)
				}
			}
			if test.expected {
				t.Logf("\t%s\tTest %d: Correctly validated the input string", success, test.testID)
			} else {
				t.Logf("\t%s\tTest %d: Correctly *invalidated* the input string", success, test.testID)
			}
		})
	}
}

func TestValidateInput(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "from hackerrank",
			input: `6
4123456789123456
5123-4567-8912-3456
61234-567-8912-3456
4123356789123456
5133-3367-8912-3456
5123 - 3567 - 8912 - 3456`,
			expected: `Valid
Valid
Invalid
Valid
Invalid
Invalid
`,
		},
	}

	for i, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("\tTest %d:\t%s", i, test.name)

			var output bytes.Buffer
			validateInput(strings.NewReader(test.input), &output)
			if output.String() != test.expected {
				t.Logf("\t\tTest %d:\tExp:\n%v", i, test.expected)
				t.Logf("\t\tTest %d:\tGot:\n%v", i, output.String())
				if diff := cmp.Diff(output.String(), test.expected); diff != "" {
					t.Logf("\t\tTest %d\tDiff:\n%v", i, diff)
				}
				t.Fatalf("\t%s\tTest %d: Incorrect Output", failed, i)
			}
			t.Logf("\t%s\tTest %d: Correct Output", success, i)
		})
	}
}
