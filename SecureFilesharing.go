package main

import (
	"fmt"

	Download "P2P-Secure-Filesharing/Libraries/download"
	JoinNetwork "P2P-Secure-Filesharing/Libraries/joinnetwork"
	PeerTypes "P2P-Secure-Filesharing/Libraries/peertypes"
	Search "P2P-Secure-Filesharing/Libraries/search"
	Upload "P2P-Secure-Filesharing/Libraries/upload"
	Peerconnect "P2P-Filesharing/Libraries/peerconnect"

)

func main() {

	fmt.Printf("Starting Secure File Sharing Application...")

	PeerTypes.PeerTypes()
	Search.Search()
	JoinNetwork.JoinNetwork("localhost")
	Download.Download()
	Upload.Upload()

}
