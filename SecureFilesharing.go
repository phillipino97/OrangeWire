package main

import (
	"flag"
	"fmt"

	Connections "P2P-Secure-Filesharing/Libraries/connections"
	Upload "P2P-Secure-Filesharing/Libraries/upload"
)

func main() {

	fmt.Printf("Starting Secure File Sharing Application...")

	firstPtr := flag.Bool("first", false, "a bool")
	addPtr := flag.String("addr", "localhost", "Address")
	portPtr := flag.String("port", "2007", "port")
	serverportPtr := flag.String("serverport", "2007", "server port")
	filepathPtr := flag.String("filepath", "2000", "filepath flag")

	flag.Parse()

	Connections.Connections(*filepathPtr)
	Upload.FILE_PATH = "./FileChunks/" + *filepathPtr

	if *firstPtr {
		Connections.RecvThread("localhost:" + *serverportPtr)
	} else {
		Connections.Init_addr = *addPtr
		Connections.Init_port = *portPtr
		go Connections.SendThread(*addPtr, *portPtr)
		Connections.RecvThread("localhost:" + *serverportPtr)
	}

}
