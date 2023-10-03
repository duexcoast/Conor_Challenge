package main

import (
	"bytes"

	"github.com/dlclark/regexp2"
)

const ccRegExpPCRE = `^(?!.*(\d)(-?\1){3})[456]\d{3}(-?\d{4}){3}$`

//  REGEX EXPLANATION:
//
//  No more than 4 consecutive repeated digits:
//
//  ^(?!.*		- Begin with a group: a negative lookahead matching
//				  any character zero or more times. This ensures we read the
//				  whole string in the lookahead.
//				- or more times.
//  (\d)		- Capture group containing a digit character
//  (-?\1){3}	- Optionally match a hyphen, then the exact match of the
//				  previous capture group (\d) (a backreference). Repeat this
//				  three times.
//  )			- End the negative lookahead. This means we won't match on any
//				  input that contains four consecutive repeated digits (not
//				  counting hyphens).
//
//  The rest of the conditions:
//
//	[456]		- String begins with either 4, 5 or 6.
//	\d{3}		- Followed by 3 digits 0-9.
//	(-?		- There may optionally be a single hyphen separating groups of
//				  four digits
//	\d{4})		- Group of four digits 0-9.
//	{3}$		- Repeat the following capture group (including the optional
//				  hyphen) three times, and end the match here. This ensures
//				  that the match consists of exactly 16 digits.

var rePCRE = regexp2.MustCompile(ccRegExpPCRE, regexp2.None)

// validatePCRE uses a single regular expression to validate all conditions.
// It accepts []byte but converts into []rune to accomodate package regexp2's
// handling of strings.
func validatePCRE(input []byte) bool {
	if isMatch, _ := rePCRE.MatchRunes(bytes.Runes(input)); isMatch {
		return true
	}
	return false
}
