package peertypes

import (
	"fmt"
)

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
}

func PeerTypes() {

	fmt.Println("Peer Types Go application")

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

	return pp

}

func CreateMiddlePeer(hash_two string, file_part_hash string, mp_address string, sp_addresses [2]string) MP {

	var mp MP

	mp.hash_two = hash_two
	mp.file_part_hash = file_part_hash
	mp.mp_address = mp_address
	mp.sp_addresses = sp_addresses

	return mp

}

func CreateStoragePeer(hash_two string, file_part_hash string) SP {

	var sp SP

	sp.hash_two = hash_two
	sp.file_part_hash = file_part_hash

	return sp

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
