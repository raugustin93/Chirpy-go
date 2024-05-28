package db

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chips: map[int]Chirp{},
		Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.Mux.Lock()
	defer db.Mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.Path, data, 0o600)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.Path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.Mux.RLock()
	defer db.Mux.RUnlock()

	dbStructure := DBStructure{}
	data, err := os.ReadFile(db.Path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}

	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		Path: path,
		Mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) GetEmails() ([]string, error) {
	db.Mux.RLock()
	defer db.Mux.RUnlock()

	emails := []string{}
	//
	// dbStructure := DBStructure{}
	// data, err := os.ReadFile(db.Path)
	// if errors.Is(err, os.ErrNotExist) {
	// 	return emails, err
	// }
	//
	// err = json.Unmarshal(data, &dbStructure)
	// if err != nil {
	// 	return emails, err
	// }
	//
	dbStructure, err := db.loadDB()
	if err != nil {
		return emails, err
	}

	for _, user := range dbStructure.Users {
		emails = append(emails, user.Email)
	}

	return emails, nil
}

func (db *DB) VerifyCredentials(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return User{}, err
			} else {
				return user, nil
			}
		}
	}

	return User{}, errors.New("Invalid credentials")
}
