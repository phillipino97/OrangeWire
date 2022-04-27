package connections

import (
	Download "OrangeWire/Libraries/download"
	PeerTypes "OrangeWire/Libraries/peertypes"
	Upload "OrangeWire/Libraries/upload"
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type peerstruct struct {
	sync.RWMutex
	allPeers map[*Peer]int
}

type growstruct struct {
	sync.RWMutex
	growMsgs map[string]string
}

type checkstruct struct {
	sync.Mutex
	connect_first bool
	root_peers    int
}

type uploadstruct struct {
	sync.Mutex
	addresses []string
	ports     []string
}

type downloadstruct struct {
	sync.Mutex
	addresses []string
	ports     []string
}

var grow growstruct
var peers peerstruct
var check checkstruct
var upload uploadstruct
var download downloadstruct
var cop string
var pp string
var hash_one string
var hash_two string
var kp string
var key_one string
var next_address string
var file_part_hash string
var times int
var all_file_hash [4]string

var OPEN_PORT string
var FILE_PATH string
var TIME string
var MAX_PEERS int
var Init_addr string
var Init_port string

var RecentSearch map[string]string

type Peer struct {
	// incoming chan string
	id         int
	incoming   string
	reader     *bufio.Reader
	writer     *bufio.Writer
	conn       net.Conn
	connection *Peer
}

func (peer *Peer) handleMessages(line string) bool {
	if strings.Contains(line, "*_+_*") {
		split := strings.Split(line, "*_+_*")
		if split[0] == "send*_=_*struct" {

			decoder := gob.NewDecoder(peer.conn)
			gen := &PeerTypes.Generic{}
			decoder.Decode(gen)
			fmt.Println("\n" + gen.ToString())
			PeerTypes.ConvertFromGeneric(*gen)
			return false

		} else if split[0] == "send*_=_*file" {

			var err error

			filename := split[1]
			filesize, err := strconv.ParseInt(split[2], 10, 64)

			const BUFFER_SIZE = 8192
			var currentbyte int64 = 0

			filebuffer := make([]byte, BUFFER_SIZE)

			file, err := os.Create(FILE_PATH + "/upload/" + strings.TrimSpace(filename))
			if err != nil {

			} else {
				for err == nil || err != io.EOF {

					peer.conn.Read(filebuffer)

					_, err = file.WriteAt(filebuffer, currentbyte)
					currentbyte += BUFFER_SIZE

					if currentbyte >= filesize {
						break
					}

				}
				file.Close()
				err := os.Truncate(FILE_PATH+"/upload/"+filename, filesize)
				if err != nil {
					log.Fatal(err)
				}

			}

			fmt.Println("\nFile Got!")

		} else if split[0] == "req*_=_*file" && strings.Contains(line, "*_+_*stream*_+_*") {

			if PeerTypes.CheckStorageStored(split[2], split[3]) != nil {

				sp := PeerTypes.CheckStorageStored(split[2], split[3])
				peer.streamFileData(sp.GetData("filename"), nil)
				temp_line, _ := peer.reader.ReadString('\n')
				temp_line = strings.Replace(temp_line, "\n", "", -1)
				//fmt.Println(temp_line)
				return false

			} else if PeerTypes.CheckMiddleStored(split[2], split[3]) != nil {

				mp := PeerTypes.CheckMiddleStored(split[2], split[3])
				temp_peer := CreateSendThread(mp.GetData("sp1"), false)
				temp_peer.sendMsg(line)
				temp_line, _ := temp_peer.reader.ReadString('\n')
				temp_line = strings.Replace(temp_line, "\n", "", -1)
				peer.streamFileData(temp_line, temp_peer)
				temp_line, _ = peer.reader.ReadString('\n')
				temp_line = strings.Replace(temp_line, "\n", "", -1)
				//fmt.Println(temp_line)
				return false

			} else if PeerTypes.CheckProxyStored(split[2]) != nil {

				pp := PeerTypes.CheckProxyStored(split[2])
				temp_peer := CreateSendThread(pp.GetData("mp1"), false)
				temp_peer.sendMsg(line)
				temp_line, _ := temp_peer.reader.ReadString('\n')
				temp_line = strings.Replace(temp_line, "\n", "", -1)
				peer.streamFileData(temp_line, temp_peer)
				temp_line, _ = peer.reader.ReadString('\n')
				temp_line = strings.Replace(temp_line, "\n", "", -1)
				//fmt.Println(temp_line)
				return false

			}

		} else if split[0] == "grow" {

			if strings.Split(split[1], "*_=_*")[0] == "search" && TIME != strings.Split(split[3], "*_=_*")[1] {

				var match2 = regexp.MustCompile("prop\\*_\\=_\\*\\d\\*_\\+_\\*$")
				check_line := match2.ReplaceAllString(line, "")
				check_line = strings.TrimSpace(check_line)

				if strings.Split(strings.Split(split[2], "*_=_*")[1], ":")[0] == "" {

					var match = regexp.MustCompile(":\\d+$")
					addr := match.ReplaceAllString(peer.conn.RemoteAddr().String(), "")
					line = strings.Replace(line, "addr*_=_*:", "addr*_=_*"+addr, 1)
					fmt.Println(line)

					check_line = match2.ReplaceAllString(line, "")
					check_line = strings.TrimSpace(check_line)

					if grow.growMsgs[strings.Split(line, "*_+_*")[2]] != check_line {

						fmt.Println("Search for file: " + strings.TrimSpace(strings.Split(split[1], "*_=_*")[1]))

						fmt.Println("Connected to source of message")

						if PeerTypes.CheckFilenameStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])) != nil {

							peer.sendMsg("found*_+_*filename*_=_*" + PeerTypes.CheckFilenameStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("title") + "*_+_*port*_=_*" + OPEN_PORT + "*_+_*")
							fmt.Println("I have the file")

						} else {

							fmt.Println("Don't have file")

							prop, _ := strconv.Atoi(strings.Split(split[4], "*_=_*")[1])

							if prop < 8 {

								grow.Lock()
								grow.growMsgs[strings.Split(line, "*_+_*")[2]] = check_line
								grow.Unlock()
								temp_line := match2.ReplaceAllString(line, "prop*_=_*"+strconv.Itoa(prop+1)+"*_+_*")
								peer.sendMsgAllNoOriginal(peer.id, temp_line)

							}

						}

					}

				} else if grow.growMsgs[strings.Split(line, "*_+_*")[2]] != check_line {

					fmt.Println("Search for file: " + strings.TrimSpace(strings.Split(split[1], "*_=_*")[1]))

					fmt.Println("Not connected to source of message")

					if PeerTypes.CheckFilenameStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])) != nil {

						fmt.Println("I have the file")

						temp_peer := CreateSendThread(strings.Split(split[2], "*_=_*")[1], false)
						temp_peer.sendMsg("found*_+_*filename*_=_*" + PeerTypes.CheckFilenameStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("title") + "*_+_*port*_=_*" + OPEN_PORT + "*_+_*")
						temp_peer.conn.Close()
						temp_peer = nil

					} else {

						fmt.Println("Don't have file")

						prop, _ := strconv.Atoi(strings.Split(split[4], "*_=_*")[1])

						if prop < 8 {

							grow.Lock()
							grow.growMsgs[strings.Split(line, "*_+_*")[2]] = check_line
							grow.Unlock()
							temp_line := match2.ReplaceAllString(line, "prop*_=_*"+strconv.Itoa(prop+1)+"*_+_*")
							peer.sendMsgAllNoOriginal(peer.id, temp_line)

						}

					}

				}

			} else if strings.Split(split[1], "*_=_*")[0] == "join" && strings.Split(split[1], "*_=_*")[1] == "req" {

				var match2 = regexp.MustCompile("prop\\*_\\=_\\*\\d\\*_\\+_\\*")
				check_line := match2.ReplaceAllString(line, "")

				if strings.Split(strings.Split(split[2], "*_=_*")[1], ":")[0] == "" {

					var match = regexp.MustCompile(":\\d+$")
					addr := match.ReplaceAllString(peer.conn.RemoteAddr().String(), "")
					line = strings.Replace(line, "addr*_=_*", "addr*_=_*"+addr, 1)

					prop, _ := strconv.Atoi(strings.Split(split[3], "*_=_*")[1])

					grow.Lock()
					if grow.growMsgs[strings.Split(line, "*_+_*")[2]] != check_line && prop < 8 {

						if peer.id > MAX_PEERS {

							grow.growMsgs[strings.Split(line, "*_+_*")[2]] = check_line
							grow.Unlock()
							temp_line := match2.ReplaceAllString(line, "prop*_=_*"+strconv.Itoa(prop+1)+"*_+_*")
							peer.sendMsgAllNoOriginal(peer.id, temp_line)
							return false

						}
						grow.Unlock()

					}

				} else {

					prop, _ := strconv.Atoi(strings.Split(split[3], "*_=_*")[1])

					if grow.growMsgs[strings.Split(line, "*_+_*")[2]] != check_line && prop < 8 {

						if peer.id > MAX_PEERS {

							grow.Lock()
							grow.growMsgs[strings.Split(line, "*_+_*")[2]] = check_line
							grow.Unlock()
							temp_line := match2.ReplaceAllString(line, "prop*_=_*"+strconv.Itoa(prop+1)+"*_+_*")
							peer.sendMsgAllNoOriginal(peer.id, temp_line)
							return false

						} else {

							temp_peer := CreateSendThread(strings.Split(split[2], "*_=_*")[1], false)
							temp_peer.sendMsg("grow*_+_*join*_=_*accept*_+_*port*_=_*" + OPEN_PORT + "*_+_*")
							temp_peer.conn.Close()
							if temp_peer.id <= MAX_PEERS {
								check.Lock()
								check.root_peers -= 1
								check.Unlock()
							}
							peers.Lock()
							for i := range peers.allPeers {
								if i.id > temp_peer.id && i.id < MAX_PEERS+1 {
									i.id -= 1
								}
							}
							delete(peers.allPeers, temp_peer)
							peers.Unlock()
							if temp_peer.connection != nil {
								temp_peer.connection.connection = nil
							}
							//fmt.Println(strconv.Itoa(temp_peer.id) + " connection Closed!")
							temp_peer = nil

						}

					}

				}

				//"grow*_+_*join*_=_*req*_+_*addr*_=_*" + OPEN_PORT + "*_+_*prop*_=_*0*_+_*"

			} else if strings.Split(split[1], "*_=_*")[0] == "join" && strings.Split(split[1], "*_=_*")[1] == "accept" {

				if check.root_peers < 2 {

					peer_port := strings.Split(split[2], "*_=_*")[1]
					var match = regexp.MustCompile(":\\d+$")
					addr := match.ReplaceAllString(peer.conn.RemoteAddr().String(), "")
					//fmt.Println(addr + ", " + peer_port[1:])
					go SendThread(addr, peer_port[1:])
					return false

				}

			} else if strings.Split(split[1], "*_=_*")[0] == "upload" && strings.Split(split[1], "*_=_*")[1] == "req" {

				var match2 = regexp.MustCompile("prop\\*_\\=_\\*\\d\\*_\\+_\\*")
				check_line := match2.ReplaceAllString(line, "")

				if strings.Split(strings.Split(split[2], "*_=_*")[1], ":")[0] == "" {

					var match = regexp.MustCompile(":\\d+$")
					addr := match.ReplaceAllString(peer.conn.RemoteAddr().String(), "")
					line = strings.Replace(line, "addr*_=_*", "addr*_=_*"+addr, 1)

					prop, _ := strconv.Atoi(strings.Split(split[3], "*_=_*")[1])

					if grow.growMsgs[strings.Split(line, "*_+_*")[2]] != check_line && prop < 5 {

						if !PeerTypes.CheckAll(strings.Split(split[5], "*_=_*")[1]) {

							peer.sendMsg("grow*_+_*upload*_=_*accept*_+_*port*_=_*" + OPEN_PORT + "*_+_*")

							grow.Lock()
							grow.growMsgs[strings.Split(line, "*_+_*")[2]] = check_line
							grow.Unlock()
							temp_line := match2.ReplaceAllString(line, "prop*_=_*"+strconv.Itoa(prop+1)+"*_+_*")
							peer.sendMsgAllNoOriginal(peer.id, temp_line)

						} else {

							grow.Lock()
							grow.growMsgs[strings.Split(line, "*_+_*")[2]] = check_line
							grow.Unlock()
							temp_line := match2.ReplaceAllString(line, "prop*_=_*"+strconv.Itoa(prop+1)+"*_+_*")
							peer.sendMsgAllNoOriginal(peer.id, temp_line)

						}

					}

				} else {

					prop, _ := strconv.Atoi(strings.Split(split[3], "*_=_*")[1])

					if grow.growMsgs[strings.Split(line, "*_+_*")[2]] != check_line && prop < 5 {

						if !PeerTypes.CheckAll(strings.Split(split[5], "*_=_*")[1]) {

							temp_peer := CreateSendThread(strings.Split(split[2], "*_=_*")[1], false)
							temp_peer.sendMsg("grow*_+_*upload*_=_*accept*_+_*port*_=_*" + OPEN_PORT + "*_+_*")
							temp_line := match2.ReplaceAllString(line, "prop*_=_*"+strconv.Itoa(prop+1)+"*_+_*")
							peer.sendMsgAllNoOriginal(peer.id, temp_line)
							temp_peer.conn.Close()
							if temp_peer.id <= MAX_PEERS {
								check.Lock()
								check.root_peers -= 1
								check.Unlock()
							}
							peers.Lock()
							for i := range peers.allPeers {
								if i.id > temp_peer.id && i.id < MAX_PEERS+1 {
									i.id -= 1
								}
							}
							delete(peers.allPeers, temp_peer)
							peers.Unlock()
							if temp_peer.connection != nil {
								temp_peer.connection.connection = nil
							}
							//fmt.Println(strconv.Itoa(temp_peer.id) + " connection Closed!")
							temp_peer = nil

						} else {

							grow.Lock()
							grow.growMsgs[strings.Split(line, "*_+_*")[2]] = check_line
							grow.Unlock()
							temp_line := match2.ReplaceAllString(line, "prop*_=_*"+strconv.Itoa(prop+1)+"*_+_*")
							peer.sendMsgAllNoOriginal(peer.id, temp_line)

						}

					}

				}

				//"grow*_+_*upload*_=_*req*_+_*addr*_=_*" + OPEN_PORT + "*_+_*prop*_=_*0*_+_*time*_=_*" + time.Now().String() + "*_+_*hash_two*_=_*" + info.Hashkeyinfo.Hash_two + "*_+_*"
				//"grow*_+_*join*_=_*req*_+_*addr*_=_*" + OPEN_PORT + "*_+_*prop*_=_*0*_+_*"

			} else if strings.Split(split[1], "*_=_*")[0] == "upload" && strings.Split(split[1], "*_=_*")[1] == "accept" {

				peer_port := strings.Split(split[2], "*_=_*")[1]
				var match = regexp.MustCompile(":\\d+$")
				addr := match.ReplaceAllString(peer.conn.RemoteAddr().String(), "")

				upload.Lock()
				upload.addresses = append(upload.addresses, addr)
				upload.ports = append(upload.ports, peer_port[1:])
				upload.Unlock()

			}

		} else if split[0] == "found" {

			peer_port := strings.Split(split[2], "*_=_*")[1]
			var match = regexp.MustCompile(":\\d+$")
			addr := match.ReplaceAllString(peer.conn.RemoteAddr().String(), "")

			download.Lock()
			download.addresses = append(download.addresses, addr)
			download.ports = append(download.ports, peer_port[1:])
			download.Unlock()

			return false
		} else if split[0] == "download" {

			if strings.Split(split[1], "*_=_*")[0] == "0" {
				hash_one := PeerTypes.CheckFilenameStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("hash_one")
				hash_two := PeerTypes.CheckFilenameStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("hash_two")
				peer.sendMsg(hash_one + "*_+_*" + hash_two)
				dp_addr := PeerTypes.CheckFilenameStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("dp1")
				peer_port := strings.Split(split[2], "*_=_*")[1]
				var match = regexp.MustCompile(":\\d+$")
				addr := match.ReplaceAllString(peer.conn.RemoteAddr().String(), "")
				temp_peer := CreateSendThread(dp_addr, false)
				temp_peer.sendMsg("download*_+_*1*_=_*" + hash_two + "*_+_*addr*_=_*" + addr + ":" + peer_port + "*_+_*")
				temp_peer.conn.Close()
				if temp_peer.id <= MAX_PEERS {
					check.Lock()
					check.root_peers -= 1
					check.Unlock()
				}
				peers.Lock()
				for i := range peers.allPeers {
					if i.id > temp_peer.id && i.id < MAX_PEERS+1 {
						i.id -= 1
					}
				}
				delete(peers.allPeers, temp_peer)
				peers.Unlock()
				if temp_peer.connection != nil {
					temp_peer.connection.connection = nil
				}
				//fmt.Println(strconv.Itoa(temp_peer.id) + " connection Closed!")
				temp_peer = nil

				return false

			} else if strings.Split(split[1], "*_=_*")[0] == "1" {
				cop_addr := PeerTypes.CheckDirectoryStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("cop1")
				pp_addr := PeerTypes.CheckDirectoryStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("pp_one1")
				temp_peer := CreateSendThread(strings.Replace(strings.TrimSpace(strings.Split(split[2], "*_=_*")[1]), "::", ":", 1), false)
				temp_peer.sendMsg("download*_=_*send*_+_*cop*_=_*" + cop_addr + "*_+_*pp*_=_*" + pp_addr + "*_+_*")
				temp_peer.conn.Close()
				if temp_peer.id <= MAX_PEERS {
					check.Lock()
					check.root_peers -= 1
					check.Unlock()
				}
				peers.Lock()
				for i := range peers.allPeers {
					if i.id > temp_peer.id && i.id < MAX_PEERS+1 {
						i.id -= 1
					}
				}
				delete(peers.allPeers, temp_peer)
				peers.Unlock()
				if temp_peer.connection != nil {
					temp_peer.connection.connection = nil
				}
				//fmt.Println(strconv.Itoa(temp_peer.id) + " connection Closed!")
				temp_peer = nil

			} else if strings.Split(split[1], "*_=_*")[0] == "2" {
				kp_addr := PeerTypes.CheckConfirmationStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[2]), strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("kp1")
				peer.sendMsg(kp_addr)

				return false

			} else if strings.Split(split[1], "*_=_*")[0] == "3" {
				key_one := PeerTypes.CheckKeyStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("key_one")
				peer.sendMsg(key_one)

				return false

			} else if strings.Split(split[1], "*_=_*")[0] == "4" {
				next := PeerTypes.CheckProxyStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("next_pp1")
				file_part := PeerTypes.CheckProxyStored(strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])).GetData("file_part")
				peer.sendMsg(next + "*_+_*" + file_part)

			}

		} else if split[0] == "download*_=_*send" {

			cop = strings.TrimSpace(strings.Split(split[1], "*_=_*")[1])
			pp = strings.TrimSpace(strings.Split(split[2], "*_=_*")[1])

		}
		//"download*_+_*0*_=_*" + filename + "*_+_*port*_=_*" + OPEN_PORT + "*_+_*"
	} else {
		//fmt.Println(line)
	}

	return true
}

func (peer *Peer) PeerRead() {
	for {
		line, err := peer.reader.ReadString('\n')
		if err == nil {
			line = strings.Replace(line, "\n", "", -1)
			exit := peer.handleMessages(line)
			if !exit {
				break
			}
		} else {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

	}

	peer.conn.Close()
	if peer.id < MAX_PEERS+1 && peer.id > 0 {
		check.Lock()
		check.root_peers = check.root_peers - 1
		check.Unlock()
	}
	peers.Lock()
	for i := range peers.allPeers {
		if i.id > peer.id && i.id < MAX_PEERS+1 {
			i.id -= 1
		}
	}
	delete(peers.allPeers, peer)
	peers.Unlock()
	if peer.connection != nil {
		peer.connection.connection = nil
	}
	//fmt.Println(strconv.Itoa(peer.id) + " connection Closed!")
	peer = nil
}

func removeDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func uploadFile(filename string, info Upload.Returnstruct) {

	upload.ports = removeDuplicateValues(upload.ports)

	gen := PeerTypes.Generic{}

	current_part := 0

	for i := 0; i < len(upload.ports); i++ {

		if upload.ports[i] != OPEN_PORT && upload.ports[i] != "NONE" {

			if i == 0 || i == 3 || i == 6 || i == 9 {

				gen.Peer_type = 6
				gen.Hash_two = info.Hashkeyinfo.Hash_two
				gen.File_part_hash = info.Fileinfo.Part_hashes[current_part]
				gen.Filename = info.Fileinfo.Names[current_part]

				peer := CreateSendThread(upload.addresses[i]+":"+upload.ports[i], false)
				peer.sendFileData(info.Fileinfo.Names[current_part])
				time.Sleep(1 * time.Second)
				peer.sendStruct(gen)

				gen.Sp_addresses[0] = upload.addresses[i] + ":" + upload.ports[i]

				for j := 0; j < len(upload.ports); j++ {
					if upload.ports[j] == upload.ports[i] {
						upload.ports[j] = "NONE"
					}
				}

				err := os.Remove(FILE_PATH + "/upload/" + info.Fileinfo.Names[current_part])
				if err != nil {
					log.Fatal(err)
				}

			} else if i == 1 || i == 4 || i == 7 || i == 10 {

				gen.Peer_type = 5
				peer := CreateSendThread(upload.addresses[i]+":"+upload.ports[i], false)
				time.Sleep(1 * time.Second)
				peer.sendStruct(gen)

				gen.Mp_addresses[0] = upload.addresses[i] + ":" + upload.ports[i]

				for j := 0; j < len(upload.ports); j++ {
					if upload.ports[j] == upload.ports[i] {
						upload.ports[j] = "NONE"
					}
				}

			} else if i == 2 || i == 5 || i == 8 || i == 11 {

				gen.Peer_type = 4
				peer := CreateSendThread(upload.addresses[i]+":"+upload.ports[i], false)
				time.Sleep(1 * time.Second)
				peer.sendStruct(gen)

				gen.Next_pp_addresses[0] = upload.addresses[i] + ":" + upload.ports[i]
				if i == 11 {
					gen.Pp_one_addresses[0] = upload.addresses[i] + ":" + upload.ports[i]
				}

				for j := 0; j < len(upload.ports); j++ {
					if upload.ports[j] == upload.ports[i] {
						upload.ports[j] = "NONE"
					}
				}

				current_part += 1

			} else if i == 12 {

				gen.Peer_type = 3
				gen.Key_one = info.Hashkeyinfo.Key_one

				peer := CreateSendThread(upload.addresses[i]+":"+upload.ports[i], false)
				time.Sleep(1 * time.Second)
				peer.sendStruct(gen)

				gen.Kp_addresses[0] = upload.addresses[i] + ":" + upload.ports[i]

				for j := 0; j < len(upload.ports); j++ {
					if upload.ports[j] == upload.ports[i] {
						upload.ports[j] = "NONE"
					}
				}

			} else if i == 13 {

				gen.Peer_type = 2
				gen.Hash_one = info.Hashkeyinfo.Hash_one

				peer := CreateSendThread(upload.addresses[i]+":"+upload.ports[i], false)
				time.Sleep(1 * time.Second)
				peer.sendStruct(gen)

				gen.Cop_addresses[0] = upload.addresses[i] + ":" + upload.ports[i]

				for j := 0; j < len(upload.ports); j++ {
					if upload.ports[j] == upload.ports[i] {
						upload.ports[j] = "NONE"
					}
				}

			} else if i == 14 {

				gen.Peer_type = 1

				peer := CreateSendThread(upload.addresses[i]+":"+upload.ports[i], false)
				time.Sleep(1 * time.Second)
				peer.sendStruct(gen)

				gen.Dp_addresses[0] = upload.addresses[i] + ":" + upload.ports[i]

				for j := 0; j < len(upload.ports); j++ {
					if upload.ports[j] == upload.ports[i] {
						upload.ports[j] = "NONE"
					}
				}

			} else if i == 15 {

				gen.Peer_type = 0
				gen.Title = filename

				peer := CreateSendThread(upload.addresses[i]+":"+upload.ports[i], false)
				time.Sleep(1 * time.Second)
				peer.sendStruct(gen)

				for j := 0; j < len(upload.ports); j++ {
					if upload.ports[j] == upload.ports[i] {
						upload.ports[j] = "NONE"
					}
				}

			}

		}

	}

}

func downloadFile(id int, filename string) {

	times = 0
	peer := CreateSendThread(download.addresses[id]+":"+download.ports[id], false)
	peer.sendMsg("download*_+_*0*_=_*" + filename + "*_+_*port*_=_*" + OPEN_PORT + "*_+_*")
	line, _ := peer.reader.ReadString('\n')
	line = strings.Replace(line, "\n", "", -1)
	hash_one = strings.TrimSpace(strings.Split(line, "*_+_*")[0])
	hash_two = strings.TrimSpace(strings.Split(line, "*_+_*")[1])
	peer.conn.Close()
	if peer.id <= MAX_PEERS {
		check.Lock()
		check.root_peers -= 1
		check.Unlock()
	}
	peers.Lock()
	for i := range peers.allPeers {
		if i.id > peer.id && i.id < MAX_PEERS+1 {
			i.id -= 1
		}
	}
	delete(peers.allPeers, peer)
	peers.Unlock()
	if peer.connection != nil {
		peer.connection.connection = nil
	}
	peer = nil

	for cop == "" && pp == "" {

	}

	time.Sleep(1 * time.Second)

	peer = CreateSendThread(cop, false)
	peer.sendMsg("download*_+_*2*_=_*" + hash_two + "*_=_*" + hash_one + "*_+_*")
	line, _ = peer.reader.ReadString('\n')
	line = strings.Replace(line, "\n", "", -1)
	kp = strings.TrimSpace(line)
	peer.conn.Close()
	if peer.id <= MAX_PEERS {
		check.Lock()
		check.root_peers -= 1
		check.Unlock()
	}
	peers.Lock()
	for i := range peers.allPeers {
		if i.id > peer.id && i.id < MAX_PEERS+1 {
			i.id -= 1
		}
	}
	delete(peers.allPeers, peer)
	peers.Unlock()
	if peer.connection != nil {
		peer.connection.connection = nil
	}
	peer = nil

	peer = CreateSendThread(kp, false)
	peer.sendMsg("download*_+_*3*_=_*" + hash_two + "*_+_*")
	line, _ = peer.reader.ReadString('\n')
	line = strings.Replace(line, "\n", "", -1)
	key_one = strings.TrimSpace(line)
	peer.conn.Close()
	if peer.id <= MAX_PEERS {
		check.Lock()
		check.root_peers -= 1
		check.Unlock()
	}
	peers.Lock()
	for i := range peers.allPeers {
		if i.id > peer.id && i.id < MAX_PEERS+1 {
			i.id -= 1
		}
	}
	delete(peers.allPeers, peer)
	peers.Unlock()
	if peer.connection != nil {
		peer.connection.connection = nil
	}
	peer = nil

	peer = CreateSendThread(pp, false)
	peer.sendMsg("download*_+_*4*_=_*" + hash_two + "*_+_*")
	line, _ = peer.reader.ReadString('\n')
	line = strings.Replace(line, "\n", "", -1)
	next_address = strings.TrimSpace(strings.Split(line, "*_+_*")[0])
	file_part_hash = strings.TrimSpace(strings.Split(line, "*_+_*")[1])

	peer.reqFileData(hash_two, file_part_hash)

	peer.conn.Close()
	if peer.id <= MAX_PEERS {
		check.Lock()
		check.root_peers -= 1
		check.Unlock()
	}
	peers.Lock()
	for i := range peers.allPeers {
		if i.id > peer.id && i.id < MAX_PEERS+1 {
			i.id -= 1
		}
	}
	delete(peers.allPeers, peer)
	peers.Unlock()
	if peer.connection != nil {
		peer.connection.connection = nil
	}
	peer = nil

	peer = CreateSendThread(next_address, false)
	peer.sendMsg("download*_+_*4*_=_*" + hash_two + "*_+_*")
	line, _ = peer.reader.ReadString('\n')
	line = strings.Replace(line, "\n", "", -1)
	next_address = strings.TrimSpace(strings.Split(line, "*_+_*")[0])
	file_part_hash = strings.TrimSpace(strings.Split(line, "*_+_*")[1])

	peer.reqFileData(hash_two, file_part_hash)

	peer.conn.Close()
	if peer.id <= MAX_PEERS {
		check.Lock()
		check.root_peers -= 1
		check.Unlock()
	}
	peers.Lock()
	for i := range peers.allPeers {
		if i.id > peer.id && i.id < MAX_PEERS+1 {
			i.id -= 1
		}
	}
	delete(peers.allPeers, peer)
	peers.Unlock()
	if peer.connection != nil {
		peer.connection.connection = nil
	}
	peer = nil

	peer = CreateSendThread(next_address, false)
	peer.sendMsg("download*_+_*4*_=_*" + hash_two + "*_+_*")
	line, _ = peer.reader.ReadString('\n')
	line = strings.Replace(line, "\n", "", -1)
	next_address = strings.TrimSpace(strings.Split(line, "*_+_*")[0])
	file_part_hash = strings.TrimSpace(strings.Split(line, "*_+_*")[1])

	peer.reqFileData(hash_two, file_part_hash)

	peer.conn.Close()
	if peer.id <= MAX_PEERS {
		check.Lock()
		check.root_peers -= 1
		check.Unlock()
	}
	peers.Lock()
	for i := range peers.allPeers {
		if i.id > peer.id && i.id < MAX_PEERS+1 {
			i.id -= 1
		}
	}
	delete(peers.allPeers, peer)
	peers.Unlock()
	if peer.connection != nil {
		peer.connection.connection = nil
	}
	peer = nil

	peer = CreateSendThread(next_address, false)
	peer.sendMsg("download*_+_*4*_=_*" + hash_two + "*_+_*")
	line, _ = peer.reader.ReadString('\n')
	line = strings.Replace(line, "\n", "", -1)
	next_address = strings.TrimSpace(strings.Split(line, "*_+_*")[0])
	file_part_hash = strings.TrimSpace(strings.Split(line, "*_+_*")[1])

	peer.reqFileData(hash_two, file_part_hash)

	peer.conn.Close()
	if peer.id <= MAX_PEERS {
		check.Lock()
		check.root_peers -= 1
		check.Unlock()
	}
	peers.Lock()
	for i := range peers.allPeers {
		if i.id > peer.id && i.id < MAX_PEERS+1 {
			i.id -= 1
		}
	}
	delete(peers.allPeers, peer)
	peers.Unlock()
	if peer.connection != nil {
		peer.connection.connection = nil
	}
	peer = nil

	Download.FILE_PATH = FILE_PATH
	Download.Download(all_file_hash, filename)

	cop = ""
	pp = ""

}

func (peer *Peer) PeerListen() {
	go peer.PeerRead()
}

func NewPeer(connection net.Conn, run_threads bool) *Peer {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	peer := &Peer{
		id:       check.root_peers + 1,
		incoming: "",
		conn:     connection,
		reader:   reader,
		writer:   writer,
	}
	if run_threads {
		peer.PeerListen()
	}

	if peer.id < MAX_PEERS+1 {
		check.Lock()
		check.root_peers += 1
		check.Unlock()
	}

	return peer
}

func Connections(filepath string) {

	pp = ""
	cop = ""
	FILE_PATH = "./FileChunks/" + filepath
	grow.growMsgs = make(map[string]string)
	peers.allPeers = make(map[*Peer]int)
	TIME = time.Now().String()
	RecentSearch = make(map[string]string)
	check.root_peers = 0
	MAX_PEERS = 5
	check.connect_first = true
	err := os.MkdirAll(FILE_PATH+"/download", os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(FILE_PATH+"/upload", os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	go handleUser()

}

func handleUser() {

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.Contains(text, "search") {

			fmt.Print("Enter filename to search:\n> ")
			search, _ := reader.ReadString('\n')
			search = strings.Replace(search, "\n", "", -1)

			sep := PeerTypes.CreateSearchPeer(search)
			send := GenerateSendMsg(0, sep.ConvertToGeneric())
			sendMsgAll(send)
			time.Sleep(1 * time.Second)
			for i := 0; i < len(download.ports); i++ {
				fmt.Println(strconv.Itoa(i) + ") " + download.addresses[i] + ":" + download.ports[i])
			}

			fmt.Print("Enter selection:\n> ")
			text, _ = reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			id, _ := strconv.Atoi(text)
			downloadFile(id, search)

		} else if strings.TrimSpace(text) == "upload" {

			fmt.Print("Enter filename to encrypt:\n> ")
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)

			filename := text
			fmt.Print("Enter encryption password:\n> ")
			text, _ = reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			pass := text
			info := Upload.Upload(filename, pass)
			sendMsgAll("grow*_+_*upload*_=_*req*_+_*addr*_=_*" + OPEN_PORT + "*_+_*prop*_=_*0*_+_*time*_=_*" + TIME + "*_+_*hash_two*_=_*" + info.Hashkeyinfo.Hash_two + "*_+_*")
			time.Sleep(2 * time.Second)
			uploadFile(strings.TrimSpace(filename), info)

		}

	}

}

func SendThread(address string, port string) {
	//establish connection

	connection, err := net.Dial("tcp", address+":"+port)
	if err != nil {

	}
	peer := NewPeer(connection, true)
	peers.Lock()
	peers.allPeers[peer] = 1
	peers.Unlock()
	peer.sendMsg("grow*_+_*join*_=_*req*_+_*addr*_=_*" + OPEN_PORT + "*_+_*prop*_=_*0*_+_*time*_=_*" + time.Now().String() + "*_+_*")

}

func CreateSendThread(address string, multi bool) *Peer {
	//establish connection
	connection, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	return NewPeer(connection, multi)

}

func RecvThread(port string) {

	OPEN_PORT = ":" + strings.Split(port, ":")[1]
	listener, err1 := net.Listen("tcp", port)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		peer := NewPeer(conn, true)
		peers.Lock()
		peers.allPeers[peer] = 1
		peers.Unlock()
		//fmt.Println(strconv.Itoa(peer.id) + " connected")
		//time.Sleep(1 * time.Second)

	}

}

func (peer *Peer) sendMsgAllNoOriginal(id int, message string) {
	for i, _ := range peers.allPeers {
		if i.id != id && i.id < MAX_PEERS+1 {
			peers.Lock()
			i.writer.WriteString(message + "\n")
			i.writer.Flush()
			peers.Unlock()
		}
	}
}

func sendMsgAll(message string) {
	for i, _ := range peers.allPeers {
		if i.conn != nil && i.id < MAX_PEERS+1 {
			peers.Lock()
			i.writer.WriteString(message + "\n")
			i.writer.Flush()
			peers.Unlock()
		}
	}
}

func (peer *Peer) sendMsg(message string) {
	peers.Lock()
	peer.writer.WriteString(message + "\n")
	peer.writer.Flush()
	peers.Unlock()
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

	size, err := os.Stat(FILE_PATH + "/upload/" + strings.TrimSpace(filename))
	peer.sendMsg("send*_=_*file*_+_*" + filename + "*_+_*" + strconv.Itoa(int(size.Size())) + "*_+_*")

	file, err := os.Open(FILE_PATH + "/upload/" + strings.TrimSpace(filename))
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
	all_file_hash[times] = filename
	times += 1

	const BUFFER_SIZE = 8192
	var currentbyte int64 = 0

	filebuffer := make([]byte, BUFFER_SIZE)

	file, err := os.Create(FILE_PATH + "/download/" + strings.TrimSpace(filename))
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
		err := os.Truncate(FILE_PATH+"/download/"+filename, filesize)
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

		size, err := os.Stat(FILE_PATH + "/upload/" + strings.TrimSpace(filename))
		peer.sendMsg("send*_=_*file*_+_*stream*_+_*" + filename + "*_+_*" + strconv.Itoa(int(size.Size())) + "*_+_*")

		file, err := os.Open(FILE_PATH + "/upload/" + strings.TrimSpace(filename))
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

func GenerateSendMsg(function int, information PeerTypes.Generic) string {

	peer := PeerTypes.ConvertFromGeneric(information)

	if function == 0 {
		//search

		if search, ok := peer.(PeerTypes.SEP); ok {

			temp := search.GetData("title")

			var return_string = ("grow*_+_*search*_=_*" + temp + "*_+_*addr*_=_*:" + OPEN_PORT + "*_+_*time*_=_*" + TIME + "*_+_*prop*_=_*0*_+_*")

			return return_string

		}

	} else if function == 1 {
		//request

	} else if function == 2 {
		//send requested data

	}

	return ""

}
