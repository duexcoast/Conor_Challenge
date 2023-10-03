package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Success and failure markers.
const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestValidate(t *testing.T) {
	tt := validateCCTestConditions()

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Logf("\tTest %d:\t%s", test.testID, name)

			got := validate([]byte(test.input))
			if got != test.expected {
				t.Logf("\t\tTest %d:\tInput: %v", test.testID, []byte(test.input))
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

func TestValidatePCRE(t *testing.T) {
	tt := validateCCTestConditions()

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Logf("\tTest %d:\t%s", test.testID, name)

			got := validatePCRE([]byte(test.input))
			if got != test.expected {
				t.Logf("\t\tTest %d:\tInput: %v", test.testID, test.input)
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
			name:     "from hackerrank problem sample input",
			input:    "6\n4123456789123456\n5123-4567-8912-3456\n61234-567-8912-3456\n4123356789123456\n5133-3367-8912-3456\n5123 - 3567 - 8912 - 3456",
			expected: "Valid\nValid\nInvalid\nValid\nInvalid\nInvalid\n",
		},
	}

	for i, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			t.Logf("\tTest %d:\t%s", i, test.name)

			var output bytes.Buffer
			validateInput(strings.NewReader(test.input), &output, validate)
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

func TestValidateInputFromFile(t *testing.T) {
	f, err := os.Open("testdata/cc_1000.txt")
	if err != nil {
		t.Fatal("couldn't open test data")
	}
	defer f.Close()

	if err = validateInput(f, io.Discard, validate); err != nil {
		t.Fatal(err.Error())
	}
}

type validateTC struct {
	testID   int
	input    string
	expected bool
}

func validateCCTestConditions() map[string]validateTC {
	tt := map[string]validateTC{
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
		"Double hyphen": {
			testID:   9,
			input:    "5491-0987--2578-3904",
			expected: false,
		},
		"Incorrect delimiter": {
			testID:   10,
			input:    "5491_0987_2578_3904",
			expected: false,
		},
		"Contains non-digit characters": {
			testID:   11,
			input:    "6O438I4033277914",
			expected: false,
		},
		"Starting character that is not [456]": {
			testID:   12,
			input:    "3081-2855-9842-8003",
			expected: false,
		},
	}

	return tt
}

// Using stdlib regexp package, and an additional loop.
// Processes input as []byte instead of string. This pays dividends when
// streaming through input using Scanner, because scanner.String() allocates on
// every call while scanner.Bytes() does not.
func BenchmarkValidateRE2(b *testing.B) {
	input := []byte("4563-6093-5432-5444")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		validate(input)
	}
}

// Using the regexp2 package, which is fully PCRE compatable. Uses a single
// regular expression to validate the data completely
func BenchmarkValidatePCRE(b *testing.B) {
	input := []byte("4563-6093-5432-5444")

	for n := 0; n < b.N; n++ {
		validatePCRE(input)
	}

}

// Stdlib regex package reading input from a file
func BenchmarkValidateInputRE2(b *testing.B) {

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		f, err := os.Open("testdata/cc_1000.txt")
		if err != nil {
			b.Fatal(err.Error())
		}
		defer f.Close()

		b.StartTimer()
		err = validateInput(f, io.Discard, validate)
		if err != nil {
			b.Fatalf(err.Error())
		}
	}
}

// regexp2 package reading from a file
func BenchmarkValidateInputPCRE(b *testing.B) {

	for n := 0; n < b.N; n++ {
		b.StopTimer()
		f, err := os.Open("testdata/cc_1000.txt")
		if err != nil {
			b.Fatal(err.Error())
		}
		defer f.Close()

		b.StartTimer()
		err = validateInput(f, io.Discard, validatePCRE)
		if err != nil {
			b.Fatalf(err.Error())
		}
	}
}
