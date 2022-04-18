package connections

import (
	PeerTypes "P2P-Secure-Filesharing/Libraries/peertypes"
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var allClients map[*Client]int
var allPeers map[*Peer]int

type Client struct {
	// incoming chan string
	id            int
	outgoing      chan string
	outgoing_this string
	reader        *bufio.Reader
	writer        *bufio.Writer
	conn          net.Conn
	connection    *Client
}

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

func (client *Client) ClientRead() {
	for {
		line, err := client.reader.ReadString('\n')
		if err == nil {
			if client.connection != nil {
				split := strings.Split(line, "*_+_*")
				if split[0] == "send*_=_*struct" {
					decoder := gob.NewDecoder(client.conn)
					gen := &PeerTypes.Generic{}
					decoder.Decode(gen)
					fmt.Println(gen.ToString())
				} else {
					client.connection.outgoing <- line
				}
			}
		} else {
			fmt.Println(err)
			break
		}

	}

	client.conn.Close()
	delete(allClients, client)
	if client.connection != nil {
		client.connection.connection = nil
	}
	client = nil
}

func (client *Client) ClientWrite() {
	for data := range client.outgoing {
		client.writer.WriteString(data)
		data = strings.ReplaceAll(data, "\n", "")
		fmt.Print(data)
		fmt.Print("\n")
		client.writer.Flush()
	}
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
			if peer.connection != nil {
				peer.connection.incoming = line
			}
			line = strings.ReplaceAll(line, "\n", "")
			fmt.Print(line)
			fmt.Print("\n")
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
}

func (client *Client) ClientListen() {
	go client.ClientRead()
	go client.ClientWrite()
}

func (peer *Peer) PeerListen() {
	go peer.PeerRead()
	go peer.PeerWrite()
}

func NewClient(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		// incoming: make(chan string),
		id:            len(allClients) + 1,
		outgoing:      make(chan string),
		outgoing_this: "",
		conn:          connection,
		reader:        reader,
		writer:        writer,
	}
	client.ClientListen()

	return client
}

func NewPeer(connection net.Conn) *Peer {
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
	peer.PeerListen()

	return peer
}

func Connections() {

	fmt.Println("Opening connections...")

}

func SendThread(address string, port string) {
	//establish connection
	connection, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		panic(err)
	}

	peer := NewPeer(connection)

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)

		help := PeerTypes.CreateSearchPeer("hello")
		newone := help.ConvertToGeneric()

		peer.sendStruct(newone)

	}

}

func RecvThread(port string) {
	allClients = make(map[*Client]int)
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
		client := NewClient(conn)
		for clientList := range allClients {
			if clientList.connection == nil {
				client.connection = clientList
				clientList.connection = client
				fmt.Println("Connected")
			}
		}
		allClients[client] = 1
		fmt.Println(len(allClients))
		client.sendMsg(client.id, "Welcome to the party!")
		client.sendMsgAll("Welcome peer " + strconv.Itoa(client.id) + " to the cluser!")
	}
}

func (client *Client) sendMsgAll(message string) {
	if client.connection != nil {
		client.connection.outgoing <- message + "\n"
	}
}

func (client *Client) sendMsg(id int, message string) {
	for i, _ := range allClients {
		if i.id == id {
			client.writer.WriteString(message + "\n")
			client.writer.Flush()
		}
	}
}

func (peer *Peer) sendMsg(message string) {

	peer.send_data = message + "\n"

}

func (client *Client) sendStruct(id int, info PeerTypes.Generic) {
	if client.connection != nil {
		client.connection.outgoing <- "send*_=_*struct*_+_*\n"
		encoder := gob.NewEncoder(client.connection.conn)
		encoder.Encode(info)
	}
}

func (peer *Peer) sendStruct(info PeerTypes.Generic) {
	peer.send_data = "send*_=_*struct*_+_*\n"
	encoder := gob.NewEncoder(peer.conn)
	encoder.Encode(info)
}

func (peer *Peer) sendFileData() {

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
