package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
)

func AESEncrypt(src string, key []byte, initialVector string) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, []byte(initialVector))
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	return crypted
}

func AESDecrypt(crypt []byte, key []byte, initialVector string) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if len(crypt) == 0 {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(initialVector))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)

	return PKCS5Trimming(decrypted)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func Hash(filename string) []byte {
	filename = strings.TrimSpace(filename)
	fileStat, err := os.Stat(filename)

	if err != nil {
		log.Fatal(err)
	}

	h1 := sha3.New512()
	h1.Write([]byte(fileStat.Name())) //hash only the file name
	sum1 := h1.Sum(nil)
	return sum1
}

func Hash2(filename string) []byte {
	filename = strings.TrimSpace(filename)
	fileStat, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	h2 := sha3.New512()
	h2.Write([]byte(fileStat.Name() + strconv.FormatInt(int64(fileStat.Size()), 10))) //hash file name plus its size
	sum2 := h2.Sum(nil)
	return sum2
}
