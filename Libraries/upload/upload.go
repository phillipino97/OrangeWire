package upload

import (
	Encryption "P2P-Secure-Filesharing/Libraries/encryption"
	"crypto/aes"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const SPLIT_SIZE = 4
const BUFFER_SIZE = aes.BlockSize
const IV = "jfielfosapjenvxz"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var FILE_PATH = ""

type splitstruct struct {
	Part_hashes [SPLIT_SIZE]string
	Names       [SPLIT_SIZE]string
}

type hashkeystruct struct {
	Hash_one string
	Hash_two string
	Key_one  string
	Key_two  string
}

type Returnstruct struct {
	Fileinfo    splitstruct
	Hashkeyinfo hashkeystruct
}

func Upload(filename string, password string) Returnstruct {

	temp := Returnstruct{}

	if len(password) < 16 {
		b := make([]byte, (32 - len(password)))
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}

		password = password + string(b)

	}

	temp.Hashkeyinfo.Key_one = password

	encryptFile(filename, password)
	temp.Fileinfo = splitFile(filename)
	temp.Hashkeyinfo.Hash_one = strings.ReplaceAll(string(Encryption.Hash(FILE_PATH+"/upload/"+strings.TrimSpace(filename))), "\n", "")
	temp.Hashkeyinfo.Hash_two = strings.ReplaceAll(string(Encryption.Hash2(FILE_PATH+"/upload/"+strings.TrimSpace(filename))), "\n", "")

	return temp

}

func encryptFile(filename string, password string) {

	new_name := strings.TrimSpace(filename) + ".enc"
	size, _ := os.Stat(FILE_PATH + "/upload/" + strings.TrimSpace(filename))

	var currentbyte int64 = 0
	filebuffer := make([]byte, BUFFER_SIZE)
	var err error

	file, err := os.Open(FILE_PATH + "/upload/" + strings.TrimSpace(filename))
	if err != nil {
		file.Close()
		return
	}

	f, err2 := os.OpenFile(FILE_PATH+"/upload/"+new_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err2 != nil {
		f.Close()
		log.Fatal(err2)
		return
	}

	for err == nil || err != io.EOF {

		n, err := file.ReadAt(filebuffer, currentbyte)

		data := Encryption.AESEncrypt(string(filebuffer[:n]), []byte(password), IV)

		_, err3 := f.Write(data)
		if err3 != nil {
			break
		}

		if err != nil || err == io.EOF {
			break
		}

		currentbyte += BUFFER_SIZE

		if currentbyte >= size.Size() {
			break
		}

	}

	f.Close()
	file.Close()

}

func splitFile(filename string) splitstruct {
	split := splitstruct{}
	filename = FILE_PATH + "/upload/" + strings.TrimSpace(filename) + ".enc"

	var currentbyte int64 = 0
	filebuffer := make([]byte, BUFFER_SIZE)
	var err error

	file, err := os.Open(filename)
	if err != nil {
		file.Close()
		return split
	}

	size, err := os.Stat(filename)
	if err != nil {
		return split
	}

	if int(size.Size())%SPLIT_SIZE == 0 {
		for i := 0; i < SPLIT_SIZE; i++ {

			size := int(size.Size()) / SPLIT_SIZE

			f, err2 := os.OpenFile(filename+strconv.Itoa(i), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err2 != nil {
				f.Close()
				log.Fatal(err2)
				return split
			}

			for err == nil || err != io.EOF {

				n, err := file.ReadAt(filebuffer, currentbyte)

				_, err3 := f.Write(filebuffer[:n])
				if err3 != nil {
					break
				}

				if err != nil || err == io.EOF {
					break
				}

				currentbyte += BUFFER_SIZE

				if int(currentbyte) >= size+(i*size) {
					break
				}

			}

			f.Close()
			b := make([]byte, 16)
			for i := range b {
				b[i] = letterBytes[rand.Intn(len(letterBytes))]
			}

			err4 := os.Rename(filename+strconv.Itoa(i), FILE_PATH+"/upload/"+string(b)+".enc"+strconv.Itoa(i))
			if err4 != nil {
				log.Fatal(err4)
			}
			split.Part_hashes[i] = strings.ReplaceAll(string(Encryption.Hash(FILE_PATH+"/upload/"+string(b)+".enc"+strconv.Itoa(i))), "\n", "")
			split.Names[i] = string(b) + ".enc" + strconv.Itoa(i)

		}

	} else {
		for i := 0; i < SPLIT_SIZE; i++ {

			extra := int(size.Size()) % SPLIT_SIZE
			size := (int(size.Size()) - extra) / SPLIT_SIZE
			final_size := size + extra

			f, err2 := os.OpenFile(filename+strconv.Itoa(i), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err2 != nil {
				f.Close()
				log.Fatal(err2)
				return split
			}

			for err == nil || err != io.EOF {

				n, err := file.ReadAt(filebuffer, currentbyte)

				_, err3 := f.Write(filebuffer[:n])
				if err3 != nil {
					break
				}

				if err != nil || err == io.EOF {
					break
				}

				currentbyte += BUFFER_SIZE

				if i == SPLIT_SIZE-1 {
					if int(currentbyte) >= final_size+(i*size) {
						break
					}
				} else {
					if int(currentbyte) >= size+(i*size) {
						break
					}
				}

			}

			f.Close()
			b := make([]byte, 16)
			for i := range b {
				b[i] = letterBytes[rand.Intn(len(letterBytes))]
			}

			err4 := os.Rename(filename+strconv.Itoa(i), FILE_PATH+"/upload/"+string(b)+".enc"+strconv.Itoa(i))
			if err4 != nil {
				log.Fatal(err4)
			}
			split.Part_hashes[i] = strings.ReplaceAll(string(Encryption.Hash(FILE_PATH+"/upload/"+string(b)+".enc"+strconv.Itoa(i))), "\n", "")
			split.Names[i] = string(b) + ".enc" + strconv.Itoa(i)

		}
	}

	file.Close()
	//err = os.Remove(filename)

	return split

}
