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
	idPtr := flag.Int("id", 0, "Peer ID")

	flag.Parse()

	PeerTypes.PeerTypes()
	Search.Search()
	JoinNetwork.JoinNetwork("localhost")
	Download.Download()
	Upload.Upload()
	Connections.Connections()

	var nothing [2]string
	var proxy [2]string
	var middle [2]string

	proxy[0] = "localhost:2008"
	middle[0] = "localhost:2007"

	if *idPtr == 1 {
		PeerTypes.CreateProxyPeer("hello", "hello1", "", nothing, proxy)
	} else if *idPtr == 2 {
		PeerTypes.CreateMiddlePeer("hello", "hello1", "", middle)
	} else if *idPtr == 3 {
		PeerTypes.CreateStoragePeer("hello", "hello1", "Avatar.mp4")
		PeerTypes.CreateFileNamePeer("Hi.jpg.copy.id", "hello", "hello1", nothing, nothing)
	}

	if *firstPtr {
		Connections.RecvThread(":" + *serverportPtr)
	} else if *serverPtr {
		go Connections.RecvThread(":" + *serverportPtr)
		Connections.SendThread(*addPtr, *portPtr)
	} else {
		Connections.SendThread(*addPtr, *portPtr)
	}

}
