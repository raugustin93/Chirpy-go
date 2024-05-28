package db

import "sync"

type DB struct {
	Path string
	Mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DBStructure struct {
	Chips map[int]Chirp `json:"chirps"`
	Users map[int]User  `json:"users"`
}
