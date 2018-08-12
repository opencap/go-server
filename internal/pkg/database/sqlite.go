package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/opencap/opencap/pkg/types"
	"golang.org/x/crypto/ed25519"
)

type SQLiteDatabase struct {
	Database
	db *sql.DB
}

func NewSQLiteDatabase(dataSource string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, err
	}

	sqldb := &SQLiteDatabase{
		db: db,
	}

	if err := sqldb.upgrade(); err != nil {
		return nil, err
	}

	return sqldb, nil
}

func (db *SQLiteDatabase) Close() error {
	return db.db.Close()
}

func (db *SQLiteDatabase) Version() (int, error) {
	var userVersion int

	if err := db.db.QueryRow(`PRAGMA user_version;`).Scan(&userVersion); err != nil {
		return 0, err
	}

	return userVersion, nil
}

func (db *SQLiteDatabase) upgrade() error {
	const newestVersion = 1

	ver, err := db.Version()
	if err != nil {
		return err
	}

	for ; ver < newestVersion; ver++ {
		var err error

		if ver == 0 {
			err = db.upgrade1()
		} else if ver == 1 {
			err = db.upgrade2()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (db *SQLiteDatabase) upgrade1() error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`CREATE TABLE addresses (
		domain   VARCHAR(255) NOT NULL,
		username VARCHAR(64)  NOT NULL,
		type     INT          NOT NULL,
		subtype  INT          NOT NULL,
		address  BLOB         NOT NULL,
		PRIMARY KEY (domain, username, type)
	);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE domains (
		domain    VARCHAR(255) NOT NULL,
		publicKey BLOB         DEFAULT NULL,
		PRIMARY KEY (domain)
	);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE users (
		domain   VARCHAR(255) NOT NULL,
		username VARCHAR(64)  NOT NULL,
		password VARCHAR(60)  NOT NULL,
		PRIMARY KEY (domain, username)
	);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`PRAGMA user_version = 1;`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *SQLiteDatabase) upgrade2() error {
	// noop - implement on db version upgrade
	return nil
}

func (db *SQLiteDatabase) GetAddress(domain, user string, typeId types.TypeId) (*Address, error) {
	var (
		subType uint8
		address []byte
	)

	row := db.db.QueryRow(`SELECT subtype, address FROM addresses WHERE domain = ? AND username = ? AND type = ?;`, domain, user, typeId)
	if err := row.Scan(&subType, &address); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &Address{
		SubTypeId:   subType,
		AddressData: address,
	}, nil
}

func (db *SQLiteDatabase) SetAddress(domain, user string, typeId types.TypeId, address *Address) error {
	res, err := db.db.Exec(`UPDATE addresses SET subtype = ? AND address = ? WHERE domain = ? AND username = ? AND type = ?;`, address.SubTypeId, address.AddressData, domain, user, typeId)
	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n > 0 {
		return nil
	}

	_, err = db.db.Exec(`INSERT INTO addresses (domain, username, type, subtype, address) VALUES (?, ?, ?, ?, ?);`, domain, user, typeId, address.SubTypeId, address.AddressData)
	return err
}

func (db *SQLiteDatabase) DeleteAddress(domain, user string, typeId types.TypeId) error {
	_, err := db.db.Exec(`DELETE FROM addresses WHERE domain = ? AND username = ? AND type = ?;`, domain, user, typeId)
	return err
}

func (db *SQLiteDatabase) GetPublicKey(domain string) (ed25519.PublicKey, error) {
	var publicKey []byte

	row := db.db.QueryRow(`SELECT publicKey FROM domains WHERE domain = ?;`, domain)
	if err := row.Scan(&publicKey); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return ed25519.PublicKey(publicKey), nil
}

func (db *SQLiteDatabase) SetPublicKey(domain string, kp ed25519.PublicKey) error {
	res, err := db.db.Exec(`UPDATE domains SET publicKey = ? WHERE domain = ?;`, domain, []byte(kp))
	if err != nil {
		return err
	}

	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n > 0 {
		return nil
	}

	_, err = db.db.Exec(`INSERT INTO domains (domain, publicKey) VALUES (?, ?);`, domain, []byte(kp))
	return err
}

func (db *SQLiteDatabase) DeletePublicKey(domain string) error {
	_, err := db.db.Exec(`DELETE FROM domains WHERE domain = ?;`, domain)
	return err
}

func (db *SQLiteDatabase) CreateUser(domain, user, hash string) error {
	_, err := db.db.Exec(`INSERT INTO users (domain, username, password) VALUES (?, ?, ?);`, domain, user, hash)
	return err
}

func (db *SQLiteDatabase) DeleteUser(domain, user string) error {
	_, err := db.db.Exec(`DELETE FROM users WHERE domain = ? AND username = ?;`, domain, user)
	return err
}

func (db *SQLiteDatabase) GetUserPassword(domain, user string) (string, error) {
	var password string

	row := db.db.QueryRow(`SELECT password FROM users WHERE domain = ? AND username = ?;`, domain, user)
	if err := row.Scan(&password); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return password, nil
}
