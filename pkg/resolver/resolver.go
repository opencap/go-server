package resolver

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"net"
	"regexp"
)

const (
	service = "opencap"
	proto   = "tcp"

	keyPublicKey = "opencap_key"
	keyDNSSig    = "opencap_dnssig"
)

var txtRegex = regexp.MustCompile("\\s*([^=]+)=([^=\\s]+)\\s*")

func lookupSRV(name string) ([]*Server, error) {
	_, srv, err := net.LookupSRV(service, proto, name)
	if err != nil {
		return nil, fmt.Errorf("srv lookup failed: %v", err)
	}

	l := make([]*Server, len(srv))
	for i, e := range srv {
		host := e.Target
		if host[len(host)-1] == '.' {
			host = host[:len(host)-1]
		}

		l[i] = &Server{
			Host: host,
			Port: e.Port,
		}
	}

	return l, nil
}

func lookupTXT(name string) (kv map[string]string, err error) {
	var list []string

	list, err = net.LookupTXT(name)
	if err != nil {
		return
	}

	kv = make(map[string]string)

	for _, item := range list {
		matches := txtRegex.FindAllStringSubmatch(item, -1)
		for _, match := range matches {
			if len(match) == 3 {
				kv[match[1]] = match[2]
			}
		}
	}

	return
}

type Result struct {
	name      string
	Servers   []*Server
	PublicKey ed25519.PublicKey
	DNSSig    bool
}

func (res *Result) GetServer() (string, uint16) {
	if len(res.Servers) == 0 {
		return res.name, 41145
	} else {
		return res.Servers[0].Host, res.Servers[0].Port
	}
}

type Server struct {
	Host string
	Port uint16
}

func Resolve(name string) (res *Result, err error) {
	res = &Result{name: name}

	res.Servers, err = lookupSRV(name)
	if err != nil {
		err = fmt.Errorf("SRV lookup failed: %v", err)
		return
	}

	var kv map[string]string
	kv, err = lookupTXT(name)
	if err != nil {
		err = fmt.Errorf("TXT lookup failed: %v", err)
		return
	}

	if kpBase64, ok := kv[keyPublicKey]; ok {
		var kp []byte
		kp, err = hex.DecodeString(kpBase64)
		if err != nil {
			err = fmt.Errorf("failed to decode camp key: %v", err)
			return
		}

		if len(kp) != ed25519.PublicKeySize {
			err = fmt.Errorf("camp key has invalid size: expected %d byte, got %d byte", ed25519.PublicKeySize, len(kp))
			return
		}

		res.PublicKey = ed25519.PublicKey(kp)
	}

	if dnssigStr, ok := kv[keyDNSSig]; ok {
		if dnssigStr == "1" {
			res.DNSSig = true
		} else {
			res.DNSSig = false
		}
	}

	return
}
