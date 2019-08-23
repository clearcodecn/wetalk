package internal

import (
	"github.com/go-xorm/xorm"
)

type Storage struct {
	engine *xorm.Engine
}

func NewSqlite3Storage(path string) (*Storage, error) {
	s := new(Storage)
	var err error
	s.engine, err = xorm.NewEngine("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err := s.engine.Ping(); err != nil {
		return nil, err
	}
	return s, nil
}

func NewMysqlStorage(dsn string) (*Storage, error) {
	s := new(Storage)
	var err error
	s.engine, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := s.engine.Ping(); err != nil {
		return nil, err
	}
	return s, nil
}