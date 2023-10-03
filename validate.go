package challenge

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

const ccRegExp = `(^[456]\d\d\d)(?:-*(\d\d\d\d)){3}$`

// Regular Expression for matching valid credit card numbers. Conditions:
//   - Must start with a 4, 5, or 6
//   - Must contain exactly 16 digits
//   - Must only consist of digits [0-9]
//   - May have digits in groups of 4, seperated by one hyphen "-"
//   - Must NOT use any other separator like " ", or "_"
//   - Must NOT have 4 or more consecutive repeated digits (NOT MET)
//
//   NOTE: Go's stdlib regexp package does not allow for lookarounds or
//   back references. This design choice was made in order to ensure that
//   the match time is always linear in length to the input string. In
//   order to validate the final condition above with regex alone (must
//   NOT have 4 or more consecutive repeated digits), we would need to use
//   either back references or negative look aheads.
//
//   In lieu of this, we can check inputs that successfully match all other
//   conditions with a simple loop to determine if there is 4 or more
//   consecutive repeated digits.
//
//   Alternatively, we can use a PCRE regex library that allows lookaheads
//   (at the expense of guaranteed linear time), which would be safe
//   assuming we control the input and validate the length of the input
//   strings.

func validateInput(r io.Reader, w io.Writer) error {
	re := regexp.MustCompile(ccRegExp)

	// Read input from reader parameter
	scanner := bufio.NewScanner(r)

	// First line of input is an int indicating the number of lines to process
	scanner.Scan()

	// n is the amount of lines to be read from stdin
	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return err
	}

	// read n lines from StdIn and perform validation
	for i := 0; i < n; i++ {
		scanner.Scan()
		ok := validateCC(re, scanner.Text())
		switch ok {
		case true:
			fmt.Fprintln(w, "Valid")
		case false:
			fmt.Fprintln(w, "Invalid")
		}
	}
	return nil
}

func validateCC(re *regexp.Regexp, input string) bool {
	if re.MatchString(input) {
		if !containsNConsecutiveRepeatedDigits(input, 4) {
			// successful validation
			return true
		}
		// Contains 4 or more consecutive repeated digits
		return false
	}
	// Failed the regexp validation
	return false
}

// Returns true if the input string contains n consecutive digits, otherwise
// returns false. This function will ignore hyphens in the input string.
// This function is not Unicode safe, but that's ok for these purposes.
func containsNConsecutiveRepeatedDigits(input string, n int) bool {
	counter := 0
	prevElem := input[0]

	for i := 1; i < len(input); i++ {
		if input[i] == '-' {
			continue
		}
		if input[i] == prevElem {
			counter++
		}
		// counter >= n-1 means there is n consecutive repeated digits
		if counter >= n-1 {
			return true
		}
		if input[i] != prevElem {
			counter = 0
		}
		prevElem = input[i]
	}
	return false
}
