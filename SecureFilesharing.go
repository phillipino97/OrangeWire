package main

import (
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

	PeerTypes.PeerTypes()
	Search.Search()
	JoinNetwork.JoinNetwork("localhost")
	Download.Download()
	Upload.Upload()
	Connections.Connections()

}
