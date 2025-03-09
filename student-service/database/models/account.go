package models

import (
	"fmt"
	"mvp-2-spms/services/models"
	"strconv"
)

type Account struct {
	Id    uint   `gorm:"column:professor_id"`
	Login string `gorm:"column:login"`
	Hash  []byte `gorm:"column:hash"`
	Salt  string `gorm:"column:salt"`
}

func (*Account) TableName() string {
	return "user_account"
}
func (a *Account) MapToUseCaseModel() models.Account {
	return models.Account{
		Login: a.Login,
		Hash:  a.Hash,
		Salt:  a.Salt,
		Id:    fmt.Sprint(a.Id),
	}
}

func (a *Account) MapUseCaseModelToThis(acc models.Account) {
	id, _ := strconv.Atoi(acc.Id)
	a.Id = uint(id)
	a.Hash = acc.Hash
	a.Login = acc.Login
	a.Salt = acc.Salt
}
