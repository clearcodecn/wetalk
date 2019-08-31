package internal

import pb "github.com/clearcodecn/wetalk/proto"

func (s *Storage) GetUserByUsername(username string) (*pb.User, error) {
	sess := s.engine.NewSession()
	defer sess.Close()
	user := new(pb.User)
	ok, err := sess.Where("username = ?", username).Get(user)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotExists
	}
	return user, nil
}

func (s *Storage) AddUser(user *pb.User) error {
	sess := s.engine.NewSession()
	defer sess.Close()
	if _, err := sess.InsertOne(user); err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteUser(user *pb.User) error {
	sess := s.engine.NewSession()
	defer sess.Close()

	if _, err := sess.Delete(user); err != nil {
		return err
	}
	return nil
}
