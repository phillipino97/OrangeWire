package connections

import (
	PeerTypes "P2P-Secure-Filesharing/Libraries/peertypes"
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var allPeers map[*Peer]int

type Peer struct {
	// incoming chan string
	id         int
	incoming   string
	reader     *bufio.Reader
	writer     *bufio.Writer
	conn       net.Conn
	connection *Peer
	send_data  string
}

func (peer *Peer) PeerWrite() {
	for {
		if peer.send_data != "" {
			peer.writer.WriteString(peer.send_data)
			peer.writer.Flush()
			peer.send_data = ""
		}
	}
}

func (peer *Peer) PeerRead() {
	for {
		line, err := peer.reader.ReadString('\n')
		if err == nil {
			line = strings.Replace(line, "\n", "", -1)
			if strings.Contains(line, "*_+_*") {
				split := strings.Split(line, "*_+_*")
				if split[0] == "send*_=_*struct" {

					decoder := gob.NewDecoder(peer.conn)
					gen := &PeerTypes.Generic{}
					decoder.Decode(gen)
					fmt.Println(gen.ToString())

				} else if split[0] == "send*_=_*file" {

					var err error

					filename := split[1]
					filesize, err := strconv.ParseInt(split[2], 10, 64)

					const BUFFER_SIZE = 8192
					var currentbyte int64 = 0

					filebuffer := make([]byte, BUFFER_SIZE)

					file, err := os.Create(strings.TrimSpace(filename + ".copy.jpg"))
					if err != nil {

					} else {
						for err == nil || err != io.EOF {

							peer.conn.Read(filebuffer)

							_, err = file.WriteAt(filebuffer, currentbyte)
							fmt.Printf("\rReceived: " + strconv.Itoa(int(currentbyte)))
							currentbyte += BUFFER_SIZE

							if currentbyte >= filesize {
								break
							}

						}
						file.Close()
						err := os.Truncate(filename+".copy.jpg", filesize)
						if err != nil {
							log.Fatal(err)
						}
						peer.sendMsg("File Received!")
					}

				} else if split[0] == "req*_=_*file" && strings.Contains(line, "*_+_*stream*_+_*") {

					if PeerTypes.CheckStorageStored(split[2], split[3]) != nil {

						sp := PeerTypes.CheckStorageStored(split[2], split[3])
						peer.streamFileData(sp.GetData("filename"), nil)
						temp_line, _ := peer.reader.ReadString('\n')
						temp_line = strings.Replace(temp_line, "\n", "", -1)
						fmt.Println(temp_line)
						break

					} else if PeerTypes.CheckMiddleStored(split[2], split[3]) != nil {

						mp := PeerTypes.CheckMiddleStored(split[2], split[3])
						temp_peer := CreateSendThread(mp.GetData("sp1"))
						temp_peer.sendMsg(line)
						temp_line, _ := temp_peer.reader.ReadString('\n')
						temp_line = strings.Replace(temp_line, "\n", "", -1)
						peer.streamFileData(temp_line, temp_peer)
						temp_line, _ = peer.reader.ReadString('\n')
						temp_line = strings.Replace(temp_line, "\n", "", -1)
						fmt.Println(temp_line)
						break

					} else if PeerTypes.CheckProxyStored(split[2], split[3]) != nil {

						pp := PeerTypes.CheckProxyStored(split[2], split[3])
						temp_peer := CreateSendThread(pp.GetData("mp1"))
						temp_peer.sendMsg(line)
						temp_line, _ := temp_peer.reader.ReadString('\n')
						temp_line = strings.Replace(temp_line, "\n", "", -1)
						peer.streamFileData(temp_line, temp_peer)
						temp_line, _ = peer.reader.ReadString('\n')
						temp_line = strings.Replace(temp_line, "\n", "", -1)
						fmt.Println(temp_line)
						break

					}

				}
			} else {
				fmt.Println(line)
			}
		} else {
			fmt.Println(err)
			break
		}

	}

	peer.conn.Close()
	delete(allPeers, peer)
	if peer.connection != nil {
		peer.connection.connection = nil
	}
	peer = nil
	fmt.Println("Connection Closed!")
}

func (peer *Peer) PeerListen() {
	go peer.PeerRead()
	go peer.PeerWrite()
}

func NewPeer(connection net.Conn, run_threads bool) *Peer {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	peer := &Peer{
		// incoming: make(chan string),
		id:        len(allPeers) + 1,
		incoming:  "",
		conn:      connection,
		reader:    reader,
		writer:    writer,
		send_data: "",
	}
	if run_threads {
		peer.PeerListen()
	}

	return peer
}

func Connections() {

	fmt.Println("Opening connections...")

}

func SendThread(address string, port string) {
	//establish connection

	reader := bufio.NewReader(os.Stdin)

	for {

		connection, err := net.Dial("tcp", address+":"+port)
		if err != nil {
			panic(err)
		}

		fmt.Print("Text to input: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		peer := NewPeer(connection, false)

		peer.reqFileData("hello", "hello1")

		connection.Close()

	}

}

func CreateSendThread(address string) *Peer {
	//establish connection
	connection, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	return NewPeer(connection, false)

}

func RecvThread(port string) {

	allPeers = make(map[*Peer]int)
	listener, err1 := net.Listen("tcp", port)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	for {
		fmt.Println("Waiting for incoming connections...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		peer := NewPeer(conn, true)
		for peerList := range allPeers {
			if peerList.connection == nil {
				peer.connection = peerList
				peerList.connection = peer
				fmt.Println("Connected")
			}
		}
		allPeers[peer] = 1
		fmt.Println(len(allPeers))

	}

}

func (peer *Peer) sendMsgAllNoOriginal(id int, message string) {
	for i, _ := range allPeers {
		if i.id != id {
			i.writer.WriteString(message + "\n")
			i.writer.Flush()
		}
	}
}

func (peer *Peer) sendMsgAll(message string) {
	for i, _ := range allPeers {
		if i.conn != nil {
			i.writer.WriteString(message + "\n")
			i.writer.Flush()
		}
	}
}

func (peer *Peer) sendMsg(message string) {
	peer.writer.WriteString(message + "\n")
	peer.writer.Flush()
}

func (peer *Peer) sendStruct(info PeerTypes.Generic) {
	peer.sendMsg("send*_=_*struct*_+_*")
	encoder := gob.NewEncoder(peer.conn)
	encoder.Encode(info)
}

func (peer *Peer) sendFileData(filename string) {

	const BUFFER_SIZE = 8192
	var currentbyte int64 = 0
	filebuffer := make([]byte, BUFFER_SIZE)
	var err error

	size, err := os.Stat(strings.TrimSpace(filename))
	peer.sendMsg("send*_=_*file*_+_*" + filename + "*_+_*" + strconv.Itoa(int(size.Size())) + "*_+_*")

	file, err := os.Open(strings.TrimSpace(filename))
	if err != nil {
		peer.conn.Write([]byte("-1"))
		file.Close()
		return
	}

	for err == nil || err != io.EOF {

		n, err := file.ReadAt(filebuffer, currentbyte)
		peer.conn.Write(filebuffer[:n])

		if err != nil || err == io.EOF {
			break
		}

		currentbyte += BUFFER_SIZE
		fmt.Printf("\rWritten: " + strconv.Itoa(int(currentbyte)))

	}

	file.Close()

	fmt.Println()

}

func (peer *Peer) reqFileData(hash_two string, file_part_hash string) {
	peer.sendMsg("req*_=_*file*_+_*stream*_+_*" + hash_two + "*_+_*" + file_part_hash + "*_+_*")
	var err error

	info, err := peer.reader.ReadString('\n')

	info = strings.Replace(info, "\n", "", -1)
	split := strings.Split(info, "*_+_*")

	filesize, err := strconv.ParseInt(split[3], 10, 64)
	filename := split[2]

	const BUFFER_SIZE = 8192
	var currentbyte int64 = 0

	filebuffer := make([]byte, BUFFER_SIZE)

	file, err := os.Create(strings.TrimSpace(filename + ".copy.jpg"))
	if err != nil {

	} else {
		for err == nil || err != io.EOF {

			peer.conn.Read(filebuffer)

			_, err = file.WriteAt(filebuffer, currentbyte)
			fmt.Printf("\rReceived: " + strconv.Itoa(int(currentbyte)))
			currentbyte += BUFFER_SIZE

			if currentbyte >= filesize {
				break
			}

		}
		file.Close()
		err := os.Truncate(filename+".copy.jpg", filesize)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println()
	peer.sendMsg("File Received!")
}

func (peer *Peer) streamFileData(filename string, temp_peer *Peer) {

	if temp_peer == nil {

		const BUFFER_SIZE = 8192
		var currentbyte int64 = 0
		filebuffer := make([]byte, BUFFER_SIZE)
		var err error

		size, err := os.Stat(strings.TrimSpace(filename))
		peer.sendMsg("send*_=_*file*_+_*stream*_+_*" + filename + "*_+_*" + strconv.Itoa(int(size.Size())) + "*_+_*")

		file, err := os.Open(strings.TrimSpace(filename))
		if err != nil {
			peer.conn.Write([]byte("-1"))
			file.Close()
			return
		}

		for err == nil || err != io.EOF {

			n, err := file.ReadAt(filebuffer, currentbyte)
			peer.conn.Write(filebuffer[:n])

			if err != nil || err == io.EOF {
				break
			}

			currentbyte += BUFFER_SIZE
			fmt.Printf("\rWritten: " + strconv.Itoa(int(currentbyte)))

		}

		file.Close()

		fmt.Println()

	} else {

		peer.sendMsg(filename)

		var err error

		split := strings.Split(filename, "*_+_*")
		filesize, err := strconv.ParseInt(split[3], 10, 64)

		const BUFFER_SIZE = 8192
		var currentbyte int64 = 0

		filebuffer := make([]byte, BUFFER_SIZE)

		for err == nil || err != io.EOF {

			temp_peer.conn.Read(filebuffer)

			peer.conn.Write(filebuffer)
			fmt.Printf("\rReceived: " + strconv.Itoa(int(currentbyte)))
			currentbyte += BUFFER_SIZE

			if currentbyte >= filesize {
				break
			}

		}
		if err != nil {
			log.Fatal(err)
		}
		temp_peer.sendMsg("File Received!")
		fmt.Println()

	}

}

func GenerateSendMsg(address string, send_hash string, function int, information PeerTypes.Generic) string {

	peer := PeerTypes.ConvertFromGeneric(information)

	if function == 0 {
		//search

		if search, ok := peer.(PeerTypes.SEP); ok {

			return "search*_=_*" + search.GetData("title") + "*_+_*addr*_=_*" + address + "*_+_*hash*_=_*" + send_hash + "*_+_*"

		}

	} else if function == 1 {
		//request

	} else if function == 2 {
		//send requested data

	}

	return ""

}
