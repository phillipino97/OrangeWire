package main

import (
	"flag"
	"fmt"

	Connections "P2P-Secure-Filesharing/Libraries/connections"
	Download "P2P-Secure-Filesharing/Libraries/download"
	JoinNetwork "P2P-Secure-Filesharing/Libraries/joinnetwork"
	PeerTypes "P2P-Secure-Filesharing/Libraries/peertypes"
	Search "P2P-Secure-Filesharing/Libraries/search"
	Upload "P2P-Secure-Filesharing/Libraries/upload"
)

func main() {

	fmt.Printf("Starting Secure File Sharing Application...")

	serverPtr := flag.Bool("server", false, "a bool")
	firstPtr := flag.Bool("first", false, "a bool")
	addPtr := flag.String("addr", "localhost", "Address")
	portPtr := flag.String("port", "2007", "port")
	serverportPtr := flag.String("serverport", "2007", "server port")

	flag.Parse()

	PeerTypes.PeerTypes()
	Search.Search()
	JoinNetwork.JoinNetwork("localhost")
	Download.Download()
	Upload.Upload()
	Connections.Connections()

	if *firstPtr {
		Connections.RecvThread(":" + *serverportPtr)
	} else if *serverPtr {
		go Connections.RecvThread(":" + *serverportPtr)
		Connections.SendThread(*addPtr, *portPtr)
	} else {
		Connections.SendThread(*addPtr, *portPtr)
	}

}
