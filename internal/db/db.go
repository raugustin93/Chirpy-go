package db

import (
	"sync"
	"time"
)

type DB struct {
	Path string
	Mux  *sync.RWMutex
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type RefreshToken struct {
	Token          string    `json:"token"`
	UserId         int       `json:"userId"`
	ExpirationTime time.Time `json:"expirationTime"`
}

type DBStructure struct {
	Chips        map[int]Chirp           `json:"chirps"`
	Users        map[int]User            `json:"users"`
	RefreshToken map[string]RefreshToken `json:"refresh_tokens"`
}
