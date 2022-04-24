package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/sha3"
	//"strings"
)

func main() {
	/*
		h := sha3.New512()
		h.Write([]byte("Enter text here to be hashed"))
		sum1 := h.Sum(nil)
		fmt.Printf("hash = %x\n", sum1)
	*/
}

func hash(filename string) {
	fileStat, err := os.Stat(filename)

	if err != nil {
		log.Fatal(err)
	}

	h1 := sha3.New512()
	h1.Write([]byte(fileStat.Name())) //hash only the file name
	sum1 := h1.Sum(nil)
	fmt.Printf("hash = %x\n", sum1)
}

func hash2(filename string) {
	fileStat, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	h2 := sha3.New512()
	h2.Write([]byte(fileStat.Name() + strconv.FormatInt(int64(fileStat.Size()), 10))) //hash file name plus its size
	sum2 := h2.Sum(nil)
	fmt.Printf("hash = %x\n", sum2)
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
