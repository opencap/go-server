package database

import (
	"github.com/opencap/opencap/pkg/types"
)

type Database interface {
	GetAddress(domain, user string, typeId types.TypeId) (*Address, error)
	SetAddress(domain, user string, typeId types.TypeId, address *Address) error
	DeleteAddress(domain, user string, typeId types.TypeId) error
	SetDomainPublicKey(kp interface{}) error
	GetDomainPublicKey(domain string) (interface{}, error)
	CreateUser(domain, user, hash string) error
	DeleteUser(domain, user string) error
	GetUserPassword(domain, user string) (string, error)
}

type Address struct {
	SubTypeId   types.SubTypeId
	AddressData []byte
}
