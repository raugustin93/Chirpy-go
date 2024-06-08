package db

import "errors"

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

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chips[id]
	if !ok {
		return Chirp{}, errors.New("ID does not exist")
	}

	return chirp, nil
}

func (db *DB) CreateChirp(body string, userId int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chips) + 1
	chirp := Chirp{
		Id:       id,
		Body:     body,
		AuthorId: userId,
	}
	dbStructure.Chips[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.Chips, id)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
