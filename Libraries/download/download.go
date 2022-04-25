package download

import (
	"crypto/aes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const BUFFER_SIZE = aes.BlockSize

var FILE_PATH string

func Download(names [4]string, filename string) {

	var err error

	file, err := os.OpenFile(FILE_PATH+"/download/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		file.Close()
		return
	}

	fmt.Println(names[0] + ", " + filename)

	for i := 0; i < len(names); i++ {

		curr := names[i]

		size, _ := os.Stat(FILE_PATH + "/download/" + strings.TrimSpace(curr))

		var currentbyte int64 = 0
		filebuffer := make([]byte, BUFFER_SIZE)

		f, err2 := os.Open(FILE_PATH + "/download/" + strings.TrimSpace(curr))
		if err2 != nil {
			f.Close()
			log.Fatal(err2)
			return
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

			if currentbyte >= size.Size() {
				break
			}

		}

		f.Close()
		err5 := os.Remove(FILE_PATH + "/download/" + strings.TrimSpace(curr))
		if err5 != nil {
			log.Fatal(err5)
		}

	}

	file.Close()

}
