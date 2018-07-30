package types

type (
	TypeId    = uint16
	SubTypeId = uint8
)

var types = map[TypeId]SubTypeChecker{
	0: getSubTypeChecker(bitcoin), // Bitcoin Cash
	1: getSubTypeChecker(nano),    // Nano / Raiblocks
}

func Check(typeId TypeId, subTypeId SubTypeId, addrLen uint) bool {
	if checker, ok := types[typeId]; !ok {
		return false
	} else {
		return checker(subTypeId, addrLen)
	}
}
