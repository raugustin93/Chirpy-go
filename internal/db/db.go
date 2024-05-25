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

type DBStructure struct {
	Chips map[int]Chirp `json:"chirps"`
}
