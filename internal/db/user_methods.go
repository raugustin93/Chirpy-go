package db

import "errors"

func (db *DB) CreateUser(body, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		Id:       id,
		Email:    body,
		Password: password,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(user User) error {
	// db.Mux.Lock()
	// defer db.Mux.Lock()

	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	if !userExists(dbStructure, user) {
		return errors.New("User wasn't found")
	}

	index := getUserIndex(dbStructure, user)

	if index < 0 {
		errors.New("Id not found")
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

	return User{}, errors.New("Couldn't find user with that user id")
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
