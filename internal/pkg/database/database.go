package database

import (
	"github.com/opencap/opencap/pkg/types"
	"golang.org/x/crypto/ed25519"
)

type Database interface {
	GetAddress(domain, user string, typeId types.TypeId) (*Address, error)
	SetAddress(domain, user string, typeId types.TypeId, address *Address) error
	DeleteAddress(domain, user string, typeId types.TypeId) error

	GetPublicKey(domain string) (ed25519.PublicKey, error)
	SetPublicKey(domain string, kp ed25519.PublicKey) error
	DeletePublicKey(domain string) error

	CreateUser(domain, user, hash string) error
	DeleteUser(domain, user string) error

	GetUserPassword(domain, user string) (string, error)
	Close() error
}

type Address struct {
	SubTypeId   types.SubTypeId
	AddressData []byte
}
