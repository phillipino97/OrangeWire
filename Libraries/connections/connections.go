package connections

import (
	PeerTypes "P2P-Secure-Filesharing/Libraries/peertypes"
	"fmt"
	"net"
)

type peerData struct {
	//int
	//string
	//[]bytes
}

func Connections() {

	fmt.Println("Opening connections...")

}

func openConn() {
	//call recvThread()
	//call sendThread()
}

func sendThread() { //include ip address in parameters
	//send peerData
}

func recvThread() {
	//recv peerData
	//open new recv thread
	//this.thread.close()
	dstream, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer dstream.Close()

	for {
		con, err := dstream.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handle(con)
	}
}

func handle(con net.Conn) {

	for {

		var data string
		var err error

		//data, err := bufio.NewReader(con).ReadString("\n")

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(data)
	}

	con.Close()
}

func sendMsgAll(message string) {

}

func sendMsg(id int, message string) {

}

func sendMsgOutsidePeers(address string, message string) {

}

func sendStruct(id int, info PeerTypes.Generic) {

}

func sendStructOutsidePeers(address string, info PeerTypes.Generic) {

}

func GenerateSendMsg(address string, send_hash string, function int, information PeerTypes.Generic) string {

	peer := PeerTypes.ConvertFromGeneric(information)

	if function == 0 {
		//search

		if search, ok := peer.(PeerTypes.SEP); ok {

			return "search*_=_*" + search.GetData("title") + "*_+_*addr*_=_*" + address + "*_+_*hash*_=_*" + send_hash

		}

	} else if function == 1 {
		//request

	} else if function == 2 {
		//send requested data

	}

	return ""

}
