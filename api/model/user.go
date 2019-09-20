package model

import "time"

type User struct {
	ID        int       `json:"id" xorm:"id pk autoincr"`
	Mobile    string    `json:"mobile" xorm:"mobile"`
	Email     string    `json:"email" xorm:"email"`
	Avatar    string    `json:"avatar" xorm:"avatar"`
	Password  string    `json:"-" xorm:"password"`
	AddVerify bool      `json:"add_verify" xorm:"add_verify"`
	CreateAt  time.Time `json:"create_at" xorm:"create_at"`
	DeleteAt  time.Time `json:"delete_at" xorm:"delete_at"`
}

func (u User) TableName() string {
	return `user`
}

func (m *Model) GetUserByEmail(email string) (user *User, err error) {
	var ok bool
	ok, err = m.engine.Where("email = ?", email).Get(user)
	if !ok {
		return nil, ErrNotFound
	}
	return
}

func (m *Model) GetUserByMobile(mobile string) (user *User, err error) {
	var ok bool
	ok, err = m.engine.Where("mobile = ?", mobile).Get(user)
	if !ok {
		return nil, ErrNotFound
	}
	return
}

func (m *Model) GetUserById(id int) (user *User, err error) {
	var ok bool
	ok, err = m.engine.Where("id = ?", id).Get(user)
	if !ok {
		return nil, ErrNotFound
	}
	return
}

func (m *Model) CreateUser(user *User) error {
	_, err := m.engine.Insert(user)
	return err
}
