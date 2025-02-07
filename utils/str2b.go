// str2b prints a string contained in a file to a formatted byte array,
// useful for including in test byte slices.
// Beware minimal error checking.
// RCL 07 February 2025
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const maxLineLength int = 90

func main() {
	if len(os.Args) != 2 {
		fmt.Println("please provide a file name to convert string to bytes")
		os.Exit(1)
	}
	contents, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s := fmt.Sprintf("%v", contents)
	words := strings.Split(s, " ")

	line := ""
	printLine := func(s string) {
		line += fmt.Sprintf("%s, ", s)
		if len(line) > maxLineLength {
			fmt.Println(strings.TrimSpace(line))
			line = ""
		}
	}

	printLine("[]byte{")
	for _, w := range words {
		w := strings.Trim(w, "[]")
		printLine(w)
	}
	fmt.Printf("%s}", line)
}
