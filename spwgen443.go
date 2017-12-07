////////////////////////////////////////////////////////////////////////////////
//
//  File           : spwgen443.go
//  Description    : This is the implementaiton file for the spwgen443 password
//                   generator program.  See assignment details.
//
//  Collaborators  : Raquel Alvarez
//  Last Modified  : 12/6/2017
//

// Package statement
package main

// Imports
import (
	"fmt"
	"os"
	"math/rand"
	"strconv"
	"time"
	"github.com/pborman/getopt"
	"bufio"
	"strings"
	"regexp"
	// There will likely be several mode APIs you need
)

// Global data
var patternval string = `pattern (set of symbols defining password)

        A pattern consists of a string of characters "xxxxx",
        where the x pattern characters include:

          d - digit
          c - upper or lower case character
          l - lower case character
          u - upper case character
          w - random word from /usr/share/dict/words (or /usr/dict/words)
              note that w# will identify a word of length #, if possible
          s - special character in ~!@#$%^&*()-_=+{}[]:;/?<>,.|\

        Note: the pattern overrides other flags, e.g., -w`

// You may want to create more global variables
var letters [52]string
var letters_upper [26]string
var letters_lower [26]string
var special_chars [29]string
var digits [10]string

//
// Functions

// Up to you to decide which functions you want to add

////////////////////////////////////////////////////////////////////////////////
//
// Function

////////////////////////////////////////////////////////////////////////////////
//
// Function     : generatePasword
// Description  : This is the function to generate the password.
//
// Inputs       : length - length of password
//                pattern - pattern of the file ("" if no pattern)
//                webflag - is this a web password?
// Outputs      : 0 if successful test, -1 if failure

func generatePasword(length int8, pattern string, webflag bool) string {

	pwd := "" // Start with nothing and add code

	fmt.Printf("Information received: [%s]\n", pattern)
	fmt.Printf("Lenght [%d]\n", length)
	fmt.Printf("Webflag [%t]\n", webflag)

	// Setup character arrays
	letters = [...]string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}
	letters_lower = [...]string{"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"}
	letters_upper = [...]string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}
	special_chars = [...]string{"~","!","@","#","$","%","^","&","*","(",")","-","_","=","+","{","}","[","]",":",";","/","?","<",">",",",".","|","\\"}
	digits = [...]string{"0","1","2","3","4","5","6","7","8","9"}

	// Setup dictionary
	// Open the file
	f, _ := os.Open("/usr/share/dict/words")
	scanner := bufio.NewScanner(f)
	dict_size := 0
	dict := []string{}
	// Regular expression for valid words
	valid_word := regexp.MustCompile(`^[a-zA-Z]+$`)
	for scanner.Scan() {
		word := scanner.Text()
		// store it in an array, if no special characters are present
		if valid_word.MatchString(word) {
			append(dict, word)
		}
	}
	dict_size = len(dict)

	// Generate random password accordingly
	if pattern != "" {
		// Generate password using pattern
		for p := 0; p < len(pattern); p++ {
			//fmt.Printf("[%s]\n", string(pattern[p]))
			pat := string(pattern[p])
			pwd_char := ""
			// Setup randomness
			rand_source := rand.NewSource(time.Now().UnixNano())
			random_generator := rand.New(rand_source)

			if pat == "d" {
				// Add a random digit to the password
				rand_digit_num := random_generator.Intn(10)
				pwd_char = digits[rand_digit_num]
			} else if pat == "c" {
				// Add a random upper/lower case character to the password
				rand_letter_num := random_generator.Intn(52)
				pwd_char = letters[rand_letter_num]
			} else if pat == "l" {
				// Add a random lower case character to the password
				rand_letter_num := random_generator.Intn(26)
				pwd_char = letters_lower[rand_letter_num]
			} else if pat == "u" {
				// Add a random upper case character to the password
				rand_letter_num := random_generator.Intn(26)
				pwd_char = letters_upper[rand_letter_num]
			} else if pat == "w" {
				// Add a random word (of a part. lenght) to the password
				// Check for numbers to compose length of word
				next_pat := string(pattern[p+1])
				l, err := strconv.Atoi(next_pat)
				total_length := 0
				if err == nil {
					total_length = l
					// Length of word to be specified
					for i := p+2; i < len(pattern); i++ {
						next_pat = string(pattern[i])
						l, err := strconv.Atoi(next_pat)
						if err == nil {
							total_length = (total_length*10) + l
						} else {
							// Found last lenght digit, end loop
							i = len(pattern)
						}
					}
					// Find random dictionary word of size word_length
					// Keep generating random numbers until the word has desired length
					not_found := true
					for not_found {
						rand_word_num := random_generator.Intn(dict_size)
						if len(dict[rand_word_num]) == total_length {
							pwd_char = dict[rand_word_num]
							not_found = false
						}
					}
				} else {
					// Find random dictionary word of any size
					// Generate random number and lookup dict array
					rand_word_num := random_generator.Intn(dict_size)
					pwd_char = dict[rand_word_num]
				}

			} else if pat == "s" {
				// Add a special character to the password
				rand_spchar_num := random_generator.Intn(29)
				pwd_char = special_chars[rand_spchar_num]

			} else {
				// We skip it, assuming it's the number after "w"
			}
			pwd = pwd + pwd_char
		}

	} else if webflag {
		// Generate password with no special characters
		// Setup randomness
		rand_source := rand.NewSource(time.Now().UnixNano())
		random_generator := rand.New(rand_source)

		// Generate length number of random values, and
		// choose accordingly in between letter, digit or sp. char
		for i := 0; i < int(length); i++ {
			rand_value := random_generator.Float64()
			pwd_char := ""
			if rand_value < 0.5 {
				// Select a random letter
				rand_letter_num := random_generator.Intn(52)
				pwd_char = letters[rand_letter_num]
			} else {
				// Select a random digit
				rand_digit_num := random_generator.Intn(10)
				pwd_char = digits[rand_digit_num]
			}

			pwd = pwd + pwd_char
		}

	} else {
		// Generate random password with equal probability
		// for all character options

		// Setup randomness
		rand_source := rand.NewSource(time.Now().UnixNano())
		random_generator := rand.New(rand_source)

		// Generate length number of random values, and
		// choose accordingly in between letter, digit or sp. char
		for i := 0; i < int(length); i++ {
			rand_value := random_generator.Float64()
			pwd_char := ""
			if rand_value < 0.33 {
				// Select a random letter
				rand_letter_num := random_generator.Intn(52)
				pwd_char = letters[rand_letter_num]
			} else if (rand_value >= 0.33 && rand_value < 0.66) {
				// Select a random digit
				rand_digit_num := random_generator.Intn(10)
				pwd_char = digits[rand_digit_num]
			} else {
				//fmt.Printf("Value [%d]\n",rand_value)
				// Select a random sp. char
				rand_spchar_num := random_generator.Intn(29)
				pwd_char = special_chars[rand_spchar_num]
			}
			pwd = pwd + pwd_char
		}
	}


	// Now return the password
	return pwd
}

////////////////////////////////////////////////////////////////////////////////
//
// Function     : main
// Description  : The main function for the password generator program
//
// Inputs       : none
// Outputs      : 0 if successful test, -1 if failure

func main() {

	// Setup options for the program content
	rand.Seed(time.Now().UTC().UnixNano())
	helpflag := getopt.Bool('h', "", "help (this menu)")
	webflag := getopt.Bool('w', "", "web flag (no symbol characters, e.g., no &*...)")
	length := getopt.String('l', "", "length of password (in characters)")
	pattern := getopt.String('p', "", patternval)

	// Now parse the command line arguments
	err := getopt.Getopt(nil)
	if err != nil {
		// Handle error
		fmt.Fprintln(os.Stderr, err)
		getopt.Usage()
		os.Exit(-1)
	}

	// Get the flags
	fmt.Printf("helpflag [%t]\n", *helpflag)
	fmt.Printf("webflag [%t]\n", *webflag)
	fmt.Printf("length [%s]\n", *length)
	fmt.Printf("pattern [%s]\n", *pattern)
	// Normally, we we use getopt.Arg{#) to get the non-flag paramters

	// Safety check length parameter
	var plength int8 = 16
	if *length != "" {
		l, err := strconv.Atoi(*length)
		if err != nil {
			fmt.Printf("Bad length passed in [%s]\n", *length)
			fmt.Fprintln(os.Stderr, err)
			getopt.Usage()
			os.Exit(-1)
		}
		if l <= 0 || l > 64 {
			plength = 16
		} else {
			plength = int8(l)
		}
	}


	// Now generate the password and print it out
	pwd := generatePasword(plength, *pattern, *webflag)
	fmt.Printf("Generated password:  %s\n", pwd)

	// Return (no return code)
	return
}
