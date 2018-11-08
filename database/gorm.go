package database

import (
	"errors"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mssql"    // This is needed for gorm to know to use microsoft sql server
	_ "github.com/jinzhu/gorm/dialects/mysql"    // This is needed for gorm to know to use mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" // This is needed for gorm to know to use postgres
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // This is needed for gorm to know to use sqlite
)

// Gorm represents a connection to GORM
type Gorm struct {
	connection *gorm.DB
}

// GetGormConnection connects to the gorm database
func GetGormConnection(dbURL, dbType string) (Gorm, error) {
	db, err := gorm.Open(dbType, dbURL)
	if err != nil {
		return Gorm{}, err
	}

	err = db.DB().Ping()
	if err != nil {
		return Gorm{}, err
	}

	return Gorm{connection: db}, nil
}

// CreateTables creates the necessary tables
func (g Gorm) CreateTables(recreate bool) {
	if recreate {
		g.connection.DropTableIfExists(&User{})
		g.connection.DropTableIfExists(&Address{})
	}

	g.connection.CreateTable(&User{})
	g.connection.CreateTable(&Address{})

	g.connection.AutoMigrate(&User{})
	g.connection.AutoMigrate(&Address{})
}

// CreateUser creates a user in the database
func (g Gorm) CreateUser(user *User) error {
	if len(user.Addresses) > 0 {
		return errors.New("Created user shouldn't have any addresses")
	}

	dbc := g.connection.Create(user)
	return dbc.Error
}

// UpdateUser updates user fields
// Does not update associated addresses
func (g Gorm) UpdateUser(user User) error {
	user.Addresses = nil

	_, err := g.GetUser(user.ID)
	if err != nil {
		return errors.New("User can't be found to update")
	}

	dbc := g.connection.Save(user)
	return dbc.Error
}

// CreateOrUpdateAddress creates an address if it doesn't exist
// If it does already exist it is updated
func (g Gorm) CreateOrUpdateAddress(user *User, address Address) error {
	retrievedAddresses, err := g.GetAddresses(*user)
	if err != nil {
		return err
	}

	retrieved, err := getAddressByType(address.AddressType, retrievedAddresses)
	if err == nil {
		address.ID = retrieved.ID
		dbc := g.connection.Model(&address).Update("address", address.Address)
		return dbc.Error
	}

	address.UserID = user.ID
	dbc := g.connection.Create(&address)
	return dbc.Error
}

// DeleteUser deletes a user
func (g Gorm) DeleteUser(user User) error {
	addresses, err := g.GetAddresses(user)
	if err != nil {
		return errors.New("Can't find associated addresses. User can't be deleted. " + string(user.ID) + err.Error())
	}
	for _, v := range addresses {
		err := g.DeleteAddress(v)
		if err != nil {
			return err
		}
	}

	dbc := g.connection.Delete(&user)
	return dbc.Error
}

// DeleteAddress deletes an address
func (g Gorm) DeleteAddress(address Address) error {
	_, err := g.GetAddress(address.ID)
	if err != nil {
		return errors.New("Address " + string(address.ID) + " not found. Can't be deleted")
	}

	dbc := g.connection.Delete(&address)
	return dbc.Error
}

// GetUser returns a user given the proper id
func (g Gorm) GetUser(id uint) (User, error) {
	user := User{}
	dbc := g.connection.Where("id = ?", id).First(&user)
	if user.ID == 0 || dbc.Error != nil {
		return User{}, errors.New("User with id " + string(id) + " not found in postgres. Can't be retrieved")
	}
	return user, nil
}

// GetUserByDomainUsername returns a user given the proper username
func (g Gorm) GetUserByDomainUsername(domain, username string) (User, error) {
	user := User{}
	dbc := g.connection.Where("domain = ? and username = ?", domain, username).First(&user)
	if user.ID == 0 || dbc.Error != nil {
		return User{}, errors.New("User with id " + username + " and domain " + domain + " not found. Can't be retrieved")
	}
	return user, nil
}

// GetAddress returns an address given the proper id
func (g Gorm) GetAddress(id uint) (Address, error) {
	address := Address{}
	dbc := g.connection.Where("id = ?", id).First(&address)
	if address.ID == 0 || dbc.Error != nil {
		return Address{}, errors.New("Address with id " + string(id) + " not found. Can't be retrieved")
	}
	return address, nil
}

// GetAddressByAddressType returns an address given the proper id
func (g Gorm) GetAddressByAddressType(user User, addressType int) (Address, error) {
	retrievedAddresses, err := g.GetAddresses(user)
	if err != nil {
		return Address{}, err
	}
	return getAddressByType(addressType, retrievedAddresses)
}

// GetAddresses gets all not-deleted addresses associated with a user (user id must be provided)
func (g Gorm) GetAddresses(user User) ([]Address, error) {
	addresses := make([]Address, 0)
	if user.ID == 0 {
		return addresses, errors.New("No user id specified")
	}

	dbc := g.connection.Raw("SELECT * FROM addresses WHERE user_id = ?", user.ID).Scan(&addresses)
	return addresses, dbc.Error
}
