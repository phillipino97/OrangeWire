package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

var numpad int //keep track of padding space

func main() {
	key := "12345678901234567890123456789012"       //256 bits (32 bytes)
	iv := "1234567890123456"                        //128 bit (16 bytes)
	plaintext := "abcdefghijklmnopqrstuvwxyzABCDEF" //plaintext
	fmt.Println("Plaintext: ", plaintext)

	cipherText := fmt.Sprintf("%v", AseEncode(plaintext, key, iv, aes.BlockSize))
	fmt.Println("Encode Result:\t", cipherText)
	fmt.Println("Decode Result:\t", AseDecode(cipherText, key, iv))

}

func AseEncode(plaintext string, key string, iv string, blockSize int) string {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := Padding([]byte(plaintext), blockSize) //plaintext with padding
	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return hex.EncodeToString(ciphertext)
}

func AseDecode(cipherText string, encKey string, iv string) (decryptedString string) {
	bKey := []byte(encKey)
	bIV := []byte(iv)
	cipherTextDecoded, err := hex.DecodeString(cipherText)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(bKey)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks([]byte(cipherTextDecoded), []byte(cipherTextDecoded))

	finaltext := cipherTextDecoded

	for i := 0; i < numpad; i++ {
		if len(finaltext) > 0 {
			finaltext = finaltext[:len(finaltext)-1] //remove padding 1 at a time
		}
	}

	return string(finaltext)
}

func Padding(ciphertext []byte, blockSize int) []byte {

	padding := (blockSize - len(ciphertext)%blockSize)
	numpad = padding //keep track of padding space

	padtext := bytes.Repeat([]byte{byte(padding)}, padding) //padding value

	return append(ciphertext, padtext...) //add all the padding value to the end of plaintext
}
