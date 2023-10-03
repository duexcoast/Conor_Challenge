package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

const ccRegExp = `^[456]\d{3}(-?\d{4}){3}$`

var re = regexp.MustCompile(ccRegExp)

// REGEX EXPLANATION:
//
//	^[456]		- String begins with either 4, 5 or 6.
//	\d{3}		- Followed by 3 digits 0-9.
//	(-?			- There may optionally be a single hyphen separating groups of
//				  four digits
//	\d{4})		- Group of four digits 0-9.
//	{3}$		- Repeat the following capture group (including the optional
//				  hyphen) three times, and end the match here. This ensures
//				  that the match consists of exactly 16 digits.

type validator func([]byte) bool

func validateInput(r io.Reader, w io.Writer, v validator) error {
	// Read input from r
	scanner := bufio.NewScanner(r)

	// First line of input is an int indicating the number of lines to process
	scanner.Scan()

	// n is the amount of lines to be read from stdin
	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return err
	}

	// Read n lines from r and perform validation. Write results to w
	for i := 0; i < n; i++ {
		scanner.Scan()
		ok := v(scanner.Bytes())
		switch ok {
		case true:
			fmt.Fprintln(w, "Valid")
		case false:
			fmt.Fprintln(w, "Invalid")
		}
	}
	return nil
}

// validate returns true if the card has been successfully validated. If the
// card is invalid it will return false.
func validate(input []byte) bool {
	if re.Match(input) {
		// For all regex matches, check if there is 4 or more consecutive
		// repeated digits
		if !nConsecElems(input, 4) {
			// Successful validation
			return true
		}
		// Contains 4 or more consecutive repeated digits
		return false
	}
	// Failed the regexp validation
	return false
}

// Returns true if the input contains n consecutive repeated digits,
// otherwise returns false. This function will ignore hyphens in the input
// string. This function is not Unicode safe, but that's ok for these purposes.
func nConsecElems(input []byte, n int) bool {
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
