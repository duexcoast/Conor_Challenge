package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if err := generateTestData(1000000); err != nil {
		log.Fatal(err)
	}
}

// generates a file of test data with n random credit card numbers, some with
// hyphens some without. The first line of the test data will contain the
// number of entries in the file. The file will have the name testdata/cc_n.txt,
// where n is the provided argument. The testdata directory must already exist.
func generateTestData(n int) error {
	fileName := fmt.Sprintf("testdata/cc_%d.txt", n)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)

	// first line is the number of entries in the file.
	numLines := fmt.Sprintf("%d\n", n)
	w.WriteString(numLines)

	for i := 0; i < n; i++ {
		w.WriteString(generateCCNumber())
	}

	w.Flush()
	f.Close()
	return nil
}

// Generates a random credit card number, with a roughly 1/3 chance of the card
// being delimited by hyphens every 4 digits.
func generateCCNumber() string {
	// slice of len 16, initialized to zero values
	ccNumber := make([]string, 16)

	// initialize Rand to source random numbers
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range ccNumber {
		// assign a random number to each slot in ccNumber
		ccNumber[i] = strconv.Itoa(r.Intn(10))
	}

	cc := strings.Join(ccNumber, "") + "\n"

	// We will make 1/3 of the numbers credit card numbers have hyphens
	oneInThreeOdds := r.Intn(3)
	if oneInThreeOdds == 2 {
		groupsOfFour := make([]string, 4)
		for i := 0; i < 4; i++ {
			// beginning index of the 4 digits we want to target on each
			// iteration (i.e. 0, 4, 8, 12)
			j := i * 4
			groupsOfFour[i] = strings.Join(ccNumber[j:j+4], "")
		}
		cc = strings.Join(groupsOfFour, "-") + "\n"
	}

	return cc
}
