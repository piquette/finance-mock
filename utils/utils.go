package utils

import "fmt"

// Log writes a msg to the console.
func Log(verbose bool, format string, a ...interface{}) {
	if !verbose {
		return
	}
	format = format + "\n"
	fmt.Printf(format, a...)
}

// Contains checks if a string exists in a slice.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
