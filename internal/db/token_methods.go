package db

import "errors"

func (db *DB) InsertRefreshToken(token RefreshToken) error {
	dbStucture, err := db.loadDB()
	if err != nil {
		return err
	}

	dbStucture.RefreshToken[token.Token] = token

	err = db.writeDB(dbStucture)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetRefrehToken(token string) (RefreshToken, error) {
	dbStucture, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	refreshToken, ok := dbStucture.RefreshToken[token]
	if !ok {
		return RefreshToken{}, errors.New("Refresh token not in database")
	}

	return refreshToken, nil
}

func (db *DB) DeleteRefreshToken(token string) error {
	dbStucture, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStucture.RefreshToken, token)

	err = db.writeDB(dbStucture)
	if err != nil {
		return err
	}

	return nil
}
