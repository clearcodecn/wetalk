package model

import (
	"errors"
	"github.com/go-xorm/xorm"
)

var (
	ErrNotFound = errors.New("not found")
)

type Model struct {
	engine *xorm.Engine
}

func NewModel(driver, dsn string) (*Model, error) {
	engine, err := xorm.NewEngine(driver, dsn)
	if err != nil {
		return nil, err
	}

	m := new(Model)
	m.engine = engine

	if err := m.engine.Ping(); err != nil {
		return nil, err
	}

	if err := m.engine.Sync2(new(User), new(VerifyCode)); err != nil {
		return nil, err
	}

	return m, nil
}
