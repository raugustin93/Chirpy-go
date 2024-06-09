package db

import "errors"

func (db *DB) CreateUser(body, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		Id:          id,
		Email:       body,
		Password:    password,
		IsChirpyRed: false,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(user User) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	if !userExists(dbStructure, user) {
		return errors.New("User wasn't found")
	}

	index := getUserIndex(dbStructure, user)

	if index < 0 {
		return errors.New("id not found")
	}

	dbStructure.Users[index] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetUser(userId int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Id == userId {
			return user, nil
		}
	}

	return User{}, errors.New("couldn't find user with that user id")
}

func userExists(dbStructure DBStructure, user User) bool {
	for _, u := range dbStructure.Users {
		if u.Id == user.Id {
			return true
		}
	}
	return false
}

func getUserIndex(dbStructure DBStructure, user User) int {
	for i, u := range dbStructure.Users {
		if u.Id == user.Id {
			return i
		}
	}
	return -1
}

func (db *DB) EnableChirpyRedForUser(userId int) error {
	user, err := db.GetUser(userId)
	if err != nil {
		return err
	}
	user.IsChirpyRed = true
	err = db.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}
