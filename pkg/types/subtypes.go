package types

type SubTypeChecker func(SubTypeId, uint) bool

func getSubTypeChecker(st subType) SubTypeChecker {
	return func(typ SubTypeId, addrLen uint) bool {
		if l, ok := st[typ]; !ok {
			return false
		} else {
			return addrLen == l
		}
	}
}

type (
	subType map[SubTypeId]uint
)

var (
	bitcoin = subType{ // Bitcoin Type
		0: 20, // Public key hash (P2PKH)
		1: 20, // Script hash (P2SH)
	}

	nano = subType{ // Nano Type
		0: 32, // Public key (standard address)
	}

	bip47 = subType{ // Reusable Payment Code (BIP-47) Type
		0: 79, // Version 1 Payment Code
		1: 79, // Version 2 Payment Code
	}

	segwit = subType{ // SegWit Type
		0: 20, // Public key hash (P2WPKH)
		1: 32, // Script hash (P2WSH)
	}
)
