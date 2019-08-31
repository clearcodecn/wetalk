package internal

import (
	pb "github.com/clearcodecn/wetalk/proto"
	"sync"
)

type StaticStorage struct {
	*Storage
	users   map[string]pb.User
	usersRw sync.RWMutex
}

func NewStaticStorage(s *Storage) (*StaticStorage, error) {
	ss := new(StaticStorage)
	ss.Storage = s

	ss.users = make(map[string]pb.User)

	return ss, nil
}

func (s *StaticStorage) GetUserByUsername(username string) (*pb.User, error) {
	s.usersRw.RLock()
	if user, ok := s.users[username]; ok {
		s.usersRw.RUnlock()
		return &pb.User{
			Id:         user.Id,
			Avatar:     user.Avatar,
			Username:   user.Username,
			Password:   user.Password,
			Nickname:   user.Nickname,
			CreateDate: user.CreateDate,
		}, nil
	}

	s.usersRw.RUnlock()
	return s.Storage.GetUserByUsername(username)
}
