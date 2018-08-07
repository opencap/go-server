package resolver

import (
	"net"
	"fmt"
	"regexp"
	"golang.org/x/crypto/ed25519"
	"encoding/hex"
)

const (
	service = "opencap"
	proto = "tcp"

	keyPublicKey = "camp"
	keyDNSSig = "dnssig"
)

var txtRegex = regexp.MustCompile("\\s*([^=]+)=([^=\\s]+)\\s*")

func lookupSRV(name string) (addrs []*net.SRV, err error) {
	_, addrs, err = net.LookupSRV(service, proto, name)
	return
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
	Servers []*net.SRV
	PublicKey ed25519.PublicKey
	DNSSig bool
}

func Resolve(name string) (res *Result, err error) {
	res = &Result{}

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