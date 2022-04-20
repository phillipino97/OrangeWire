package peertypes

import (
	"fmt"
	"strconv"
)

var local_file_information []SP
var local_proxy_information []PP
var local_middle_information []MP

type Peers interface {
	GetData() string
	ConvertToGeneric() Generic
	ToString() string
}

type FNP struct {
	title         string
	hash_one      string
	hash_two      string
	dp_addresses  [2]string
	fnp_addresses [2]string
}

type DP struct {
	hash_two         string
	cop_addresses    [2]string
	pp_one_addresses [2]string
	dp_addresses     [2]string
}

type COP struct {
	hash_one      string
	hash_two      string
	cop_addresses [2]string
	kp_addresses  [2]string
}

type KP struct {
	hash_two     string
	key_one      string
	key_two      string
	kp_addresses [2]string
}

type PP struct {
	hash_two          string
	file_part_hash    string
	same_pp_address   string
	next_pp_addresses [2]string
	mp_addresses      [2]string
}

type MP struct {
	hash_two       string
	file_part_hash string
	mp_address     string
	sp_addresses   [2]string
}

type SP struct {
	hash_two       string
	file_part_hash string
	filename       string
}

type Generic struct {
	Peer_type         int
	Title             string
	Hash_one          string
	Hash_two          string
	Dp_addresses      [2]string
	Fnp_addresses     [2]string
	Cop_addresses     [2]string
	Pp_one_addresses  [2]string
	Kp_addresses      [2]string
	Key_one           string
	Key_two           string
	File_part_hash    string
	Same_pp_address   string
	Next_pp_addresses [2]string
	Mp_addresses      [2]string
	Sp_addresses      [2]string
	Mp_address        string
	Filename          string
}

type SEP struct {
	title string
}

func (gen Generic) ToString() string {
	return strconv.Itoa(gen.Peer_type) + "\n" + gen.Title
}

func (fnp FNP) GetData(data_type string) string {

	if data_type == "title" {
		return fnp.title
	} else if data_type == "hash_one" {
		return fnp.hash_one
	} else if data_type == "hash_two" {
		return fnp.hash_two
	} else if data_type == "dp1" {
		return fnp.dp_addresses[0]
	} else if data_type == "dp2" {
		return fnp.dp_addresses[1]
	} else if data_type == "fnp1" {
		return fnp.fnp_addresses[0]
	} else if data_type == "fnp2" {
		return fnp.fnp_addresses[1]
	}

	return ""

}

func (fnp FNP) ConvertToGeneric() Generic {

	var gen Generic

	gen.Peer_type = 0
	gen.Title = fnp.title
	gen.Hash_one = fnp.hash_one
	gen.Hash_two = fnp.hash_two
	gen.Dp_addresses = fnp.dp_addresses
	gen.Fnp_addresses = fnp.fnp_addresses

	return gen

}

func (dp DP) GetData(data_type string) string {

	if data_type == "hash_two" {
		return dp.hash_two
	} else if data_type == "cop1" {
		return dp.cop_addresses[0]
	} else if data_type == "cop2" {
		return dp.cop_addresses[1]
	} else if data_type == "dp1" {
		return dp.dp_addresses[0]
	} else if data_type == "dp2" {
		return dp.dp_addresses[1]
	} else if data_type == "pp_one1" {
		return dp.pp_one_addresses[1]
	} else if data_type == "pp_one2" {
		return dp.pp_one_addresses[1]
	}

	return ""

}

func (dp DP) ConvertToGeneric() Generic {

	var gen Generic

	gen.Peer_type = 1
	gen.Hash_two = dp.hash_two
	gen.Cop_addresses = dp.cop_addresses
	gen.Dp_addresses = dp.dp_addresses
	gen.Pp_one_addresses = dp.pp_one_addresses

	return gen

}

func (cop COP) GetData(data_type string) string {

	if data_type == "hash_one" {
		return cop.hash_one
	} else if data_type == "hash_two" {
		return cop.hash_two
	} else if data_type == "kp1" {
		return cop.kp_addresses[0]
	} else if data_type == "kp2" {
		return cop.kp_addresses[1]
	} else if data_type == "cop1" {
		return cop.cop_addresses[0]
	} else if data_type == "cop2" {
		return cop.cop_addresses[1]
	}

	return ""

}

func (cop COP) ConvertToGeneric() Generic {

	var gen Generic

	gen.Peer_type = 2
	gen.Hash_one = cop.hash_one
	gen.Hash_two = cop.hash_two
	gen.Kp_addresses = cop.kp_addresses
	gen.Cop_addresses = cop.cop_addresses

	return gen

}

func (kp KP) GetData(data_type string) string {

	if data_type == "hash_two" {
		return kp.hash_two
	} else if data_type == "kp1" {
		return kp.kp_addresses[0]
	} else if data_type == "kp2" {
		return kp.kp_addresses[1]
	} else if data_type == "key_one" {
		return kp.key_one
	} else if data_type == "key_two" {
		return kp.key_two
	}

	return ""

}

func (kp KP) ConvertToGeneric() Generic {

	var gen Generic

	gen.Peer_type = 3
	gen.Hash_two = kp.hash_two
	gen.Kp_addresses = kp.kp_addresses
	gen.Key_one = kp.key_one
	gen.Key_two = kp.key_two

	return gen

}

func (pp PP) GetData(data_type string) string {

	if data_type == "hash_two" {
		return pp.hash_two
	} else if data_type == "file_part" {
		return pp.file_part_hash
	} else if data_type == "same_pp" {
		return pp.same_pp_address
	} else if data_type == "next_pp1" {
		return pp.next_pp_addresses[0]
	} else if data_type == "next_pp2" {
		return pp.next_pp_addresses[1]
	} else if data_type == "mp1" {
		return pp.mp_addresses[0]
	} else if data_type == "mp2" {
		return pp.mp_addresses[1]
	}

	return ""

}

func (pp PP) ConvertToGeneric() Generic {

	var gen Generic

	gen.Peer_type = 4
	gen.Hash_two = pp.hash_two
	gen.File_part_hash = pp.file_part_hash
	gen.Same_pp_address = pp.same_pp_address
	gen.Next_pp_addresses = pp.next_pp_addresses
	gen.Mp_addresses = pp.mp_addresses

	return gen

}

func (mp MP) GetData(data_type string) string {

	if data_type == "hash_two" {
		return mp.hash_two
	} else if data_type == "file_part" {
		return mp.file_part_hash
	} else if data_type == "mp" {
		return mp.mp_address
	} else if data_type == "sp1" {
		return mp.sp_addresses[0]
	} else if data_type == "sp2" {
		return mp.sp_addresses[1]
	}

	return ""

}

func (mp MP) ConvertToGeneric() Generic {

	var gen Generic

	gen.Peer_type = 5
	gen.Hash_two = mp.hash_two
	gen.File_part_hash = mp.file_part_hash
	gen.Mp_address = mp.mp_address
	gen.Sp_addresses = mp.sp_addresses

	return gen

}

func (sp SP) GetData(data_type string) string {

	if data_type == "hash_two" {
		return sp.hash_two
	} else if data_type == "file_part" {
		return sp.file_part_hash
	} else if data_type == "filename" {
		return sp.filename
	}

	return ""

}

func (sp SP) ConvertToGeneric() Generic {

	var gen Generic

	gen.Peer_type = 6
	gen.Hash_two = sp.hash_two
	gen.File_part_hash = sp.file_part_hash
	gen.Filename = sp.filename

	return gen

}

func (sep SEP) GetData(data_type string) string {

	if data_type == "title" {
		return sep.title
	}

	return ""

}

func (sep SEP) ConvertToGeneric() Generic {

	var gen Generic

	gen.Peer_type = 7
	gen.Title = sep.title
	gen.Cop_addresses[0] = ""
	gen.Cop_addresses[1] = ""
	gen.Dp_addresses[0] = ""
	gen.Dp_addresses[1] = ""
	gen.File_part_hash = ""
	gen.Fnp_addresses[0] = ""
	gen.Fnp_addresses[1] = ""
	gen.Hash_one = ""
	gen.Hash_two = ""
	gen.Key_one = ""
	gen.Key_two = ""
	gen.Kp_addresses[0] = ""
	gen.Kp_addresses[1] = ""
	gen.Mp_address = ""
	gen.Mp_addresses[0] = ""
	gen.Mp_addresses[1] = ""
	gen.Next_pp_addresses[0] = ""
	gen.Next_pp_addresses[1] = ""
	gen.Pp_one_addresses[0] = ""
	gen.Pp_one_addresses[1] = ""
	gen.Same_pp_address = ""
	gen.Sp_addresses[0] = ""
	gen.Sp_addresses[1] = ""
	gen.Filename = ""

	return gen

}

func PeerTypes() {

	fmt.Println("Peer Types Go application")

}

func ConvertFromGeneric(peer Generic) interface{} {

	switch peer.Peer_type {

	case 0:
		return CreateFileNamePeer(peer.Title, peer.Hash_one, peer.Hash_two, peer.Dp_addresses, peer.Fnp_addresses)
	case 1:
		return CreateDirectoryPeer(peer.Hash_two, peer.Cop_addresses, peer.Pp_one_addresses, peer.Dp_addresses)
	case 2:
		return CreateConfirmationPeer(peer.Hash_one, peer.Hash_two, peer.Cop_addresses, peer.Kp_addresses)
	case 3:
		return CreateKeyPeer(peer.Hash_two, peer.Key_one, peer.Key_two, peer.Kp_addresses)
	case 4:
		return CreateProxyPeer(peer.Hash_two, peer.File_part_hash, peer.Same_pp_address, peer.Next_pp_addresses, peer.Mp_addresses)
	case 5:
		return CreateMiddlePeer(peer.Hash_two, peer.File_part_hash, peer.Mp_address, peer.Sp_addresses)
	case 6:
		return CreateStoragePeer(peer.Hash_two, peer.File_part_hash, peer.Filename)
	case 7:
		return CreateSearchPeer(peer.Title)
	}

	return nil

}

func CheckStorageStored(hash_two string, file_part_hash string) *SP {

	for _, data := range local_file_information {

		if data.hash_two == hash_two && data.file_part_hash == file_part_hash {

			temp := &data
			return temp

		}

	}

	return nil
}

func CheckProxyStored(hash_two string, file_part_hash string) *PP {

	for _, data := range local_proxy_information {

		if data.hash_two == hash_two && data.file_part_hash == file_part_hash {

			temp := &data
			return temp

		}

	}

	return nil
}

func CheckMiddleStored(hash_two string, file_part_hash string) *MP {

	for _, data := range local_middle_information {

		if data.hash_two == hash_two && data.file_part_hash == file_part_hash {

			temp := &data
			return temp

		}

	}

	return nil
}

func CreateFileNamePeer(title string, hash_one string, hash_two string, dp_addresses [2]string, fnp_addresses [2]string) FNP {

	var fnp FNP

	fnp.title = title
	fnp.hash_one = hash_one
	fnp.hash_two = hash_two
	fnp.dp_addresses = dp_addresses
	fnp.fnp_addresses = fnp_addresses

	return fnp

}

func CreateDirectoryPeer(hash_two string, cop_addresses [2]string, pp_one_addresses [2]string, dp_addresses [2]string) DP {

	var dp DP

	dp.hash_two = hash_two
	dp.cop_addresses = cop_addresses
	dp.pp_one_addresses = pp_one_addresses
	dp.dp_addresses = dp_addresses

	return dp

}

func CreateConfirmationPeer(hash_one string, hash_two string, cop_addresses [2]string, kp_addresses [2]string) COP {

	var cop COP

	cop.hash_one = hash_one
	cop.hash_two = hash_two
	cop.cop_addresses = cop_addresses
	cop.kp_addresses = kp_addresses

	return cop

}

func CreateKeyPeer(hash_two string, key_one string, key_two string, kp_addresses [2]string) KP {

	var kp KP

	kp.hash_two = hash_two
	kp.key_one = key_one
	kp.key_two = key_two
	kp.kp_addresses = kp_addresses

	return kp

}

func CreateProxyPeer(hash_two string, file_part_hash string, same_pp_address string, next_pp_addresses [2]string, mp_addresses [2]string) PP {

	var pp PP

	pp.hash_two = hash_two
	pp.file_part_hash = file_part_hash
	pp.same_pp_address = same_pp_address
	pp.next_pp_addresses = next_pp_addresses
	pp.mp_addresses = mp_addresses

	local_proxy_information = append(local_proxy_information, pp)

	return pp

}

func CreateMiddlePeer(hash_two string, file_part_hash string, mp_address string, sp_addresses [2]string) MP {

	var mp MP

	mp.hash_two = hash_two
	mp.file_part_hash = file_part_hash
	mp.mp_address = mp_address
	mp.sp_addresses = sp_addresses

	local_middle_information = append(local_middle_information, mp)

	return mp

}

func CreateStoragePeer(hash_two string, file_part_hash string, filename string) SP {

	var sp SP

	sp.hash_two = hash_two
	sp.file_part_hash = file_part_hash
	sp.filename = filename

	local_file_information = append(local_file_information, sp)

	return sp

}

func CreateSearchPeer(title string) SEP {

	var sep SEP

	sep.title = title

	return sep

}

func CreateUploadPeer() {

}

func CreateDownloadPeer() {

}

func CheckFileNamePeer(title string, fnp_test FNP) bool {

	if title == fnp_test.title {
		return true
	}

	return false

}

func CheckDirectoryPeer(hash_two string, dp_test DP) bool {

	if hash_two == dp_test.hash_two {
		return true
	}

	return false

}

func CheckConfirmationPeer(hash_one string, hash_two string, cop_test COP) bool {

	if hash_one == cop_test.hash_one && hash_two == cop_test.hash_two {
		return true
	}

	return false

}

func CheckKeyPeer(hash_two string, kp_test KP) bool {

	if hash_two == kp_test.hash_two {
		return true
	}

	return false

}

func CheckProxyPeer(hash_two string, pp_test PP) bool {

	if hash_two == pp_test.hash_two {
		return true
	}

	return false

}

func CheckMiddlePeer(hash_two string, file_part_hash string, mp_test MP) bool {

	if hash_two == mp_test.hash_two && file_part_hash == mp_test.file_part_hash {
		return true
	}

	return false

}

func CheckStoragePeer(hash_two string, file_part_hash string, sp_test SP) bool {

	if hash_two == sp_test.hash_two && file_part_hash == sp_test.file_part_hash {
		return true
	}

	return false

}

func CheckUploadPeer() {

}

func CheckDownloadPeer() {

}
