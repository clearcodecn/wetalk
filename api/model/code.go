package model

import (
	"github.com/go-xorm/xorm"
	"time"
)

type VerifyCode struct {
	Id         int       `json:"id" xorm:"id pk autoincr"`
	Code       string    `json:"code" xorm:"code"`
	Info       string    `json:"info" xorm:"info"`
	Type       int       `json:"type" xorm:"type"`
	Verified   bool      `json:"verified" xorm:"verified"`
	CreateTime time.Time `json:"create_time" xorm:"create_time"`
}

func (VerifyCode) Table() string {
	return "verify_code"
}

const (
	CodeRegister = iota
	CodeFindPassword
)

func (m *Model) CreateVerifyCode(vc *VerifyCode) error {
	_, err := m.engine.Insert(vc)
	return err
}

func (m *Model) VerifyCode(vc *VerifyCode) bool {
	_, err := m.engine.Transaction(func(session *xorm.Session) (i interface{}, e error) {
		var newVc = new(VerifyCode)
		t := time.Now().Add(-time.Minute * 30)
		i, e = session.Where("code = ? and info = ? and create_time >= ? and verified = ? and type = ?", vc.Code, vc.Info, t, false, vc.Type).Get(newVc)
		if i.(bool) {
			newVc.Verified = true
			return session.Cols("verified").Update(newVc)
		}
		return
	})
	if err != nil {
		return false
	}
	return true
}
