package storage

import (
	"bytes"
	"log"

	config "github.com/fox-one/f1db/config"
	ipfs "github.com/ipfs/go-ipfs-api"
)

var sh *ipfs.Shell

func InitIpfs() *ipfs.Shell {
	sh = ipfs.NewShell(config.GetConfig().Ipfs.URL)
	return sh
}

func WriteToIpfs(content []byte, isPin bool) (string, error) {
	var cid string
	var err error
	if isPin {
		cid, err = sh.Add(bytes.NewReader(content))
	} else {
		cid, err = sh.AddNoPin(bytes.NewReader(content))
	}
	if err != nil {
		log.Panicf("WriteToIpfs error: %s", err)
		return "", err
	}
	return cid, err
}

func ReadFromIpfs(hash string) ([]byte, error) {
	resp, err := sh.Cat(hash)
	if err != nil {
		log.Panicf("ReadFromIpfs error: %s", err)
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp)
	data := buf.Bytes()
	return data, nil
}

func PinToIpfs(path string) (string, error) {
	err := sh.Pin(path)
	if err != nil {
		log.Panicf("PinToIpfs error: %s", err)
		return "", err
	}
	return path, err
}
