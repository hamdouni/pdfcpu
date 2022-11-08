/*
Split a PDF file into single pages in same folder.

Pass the file path as a parameter (or drag and drop on the executable).

To cross-compile to Windows

	GOOS=windows GOARCH=386 go build .
*/

package main

import (
	"fmt"
	"os"

	"github.com/hamdouni/pdfcpu/pkg/api"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Drag and drop a file or give filepath in parameter!")
		waitForKey()
		os.Exit(1)
	} else if !isFile(os.Args[1]) {
		fmt.Println("File does not exist!")
		waitForKey()
		os.Exit(1)
	}

	file := os.Args[1]
	err := api.SplitFile(file, ".", 1, nil)
	if err != nil {
		fmt.Printf("Error processing file %v : receive error : %v\n", file, err)
		waitForKey()
		os.Exit(1)
	}

}
func waitForKey() {
	fmt.Print("\nPress Enter to continue...")
	var buf [1]byte
	os.Stdin.Read(buf[:])
}
func isFile(file string) bool {
	s, err := os.Stat(file)
	if err != nil || s.IsDir() {
		return false
	}
	return true
}
