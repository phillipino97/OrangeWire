package main

import (
	"fmt"

	"golang.org/x/crypto/sha3"
	/*
		"os"
		"strings"
	*/)

func main() {

	h := sha3.New512()
	h.Write([]byte("Enter text here to be hashed"))
	sum := h.Sum(nil)
	fmt.Printf("hash = %x\n", sum)
}

/*
func getInfo(filename string) {
	const BUFFER_SIZE = 8192
	var currentbyte int64 = 0
	filebuffer := make([]byte, BUFFER_SIZE)
	var err error

	size, err := os.Stat(strings.TrimSpace(filename))
}
*/
