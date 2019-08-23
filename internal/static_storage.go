package internal

import "sync"

type StaticStorage struct {
	*Storage
	users   map[string]string
	usersRw sync.Mutex
}

func NewStaticStorage(s *Storage) (*StaticStorage, error) {
	ss := new(StaticStorage)
	ss.Storage = s

	ss.users = make(map[string]string)

	return ss, nil
}