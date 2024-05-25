package db

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chips: map[int]Chirp{},
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

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chips) + 1
	chirp := Chirp{
		Id:   id,
		Body: body,
	}
	dbStructure.Chips[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chips))
	for _, chirp := range dbStructure.Chips {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}
