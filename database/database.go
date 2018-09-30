package database

import "time"

// Model overrides gorm.Model
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User of the api
type User struct {
	Model
	Username  string    `gorm:"type:varchar(30);unique_index:idx_domain_username;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Domain    string    `gorm:"not null;unique_index:idx_domain_username" json:"domain"`
	Addresses []Address `json:"addresses"`
}

// Address represents a crypto address and it's address type
type Address struct {
	Model
	UserID      uint   `gorm:"not null;unique_index:idx_userid_type"`
	Address     string `gorm:"not null" json:"address"`
	AddressType int    `gorm:"not null;unique_index:idx_userid_type" json:"address_type"`
}

// Database represents the functionality that any peristance layer for this
// server must satisfy
type Database interface {
	CreateTables(bool)
	CreateUser(*User) error
	UpdateUser(User) error
	CreateOrUpdateAddress(*User, Address) error
	DeleteUser(user User) error
	DeleteAddress(address Address) error
	GetUser(id uint) (User, error)
	GetUserByDomainUsername(domain, username string) (User, error)
	GetAddress(id uint) (Address, error)
	GetAddressByAddressType(user User, addressType int) (Address, error)
	GetAddresses(user User) ([]Address, error)
}
