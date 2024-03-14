package entities

import (
	"github.com/universalmacro/common/auth"
	"github.com/universalmacro/common/snowflake"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Account string `json:"account" gorm:"type:VARCHAR(128);uniqueIndex"`
	auth.Password
}

var accountIdGenerator = snowflake.NewIdGenertor(0)

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.Model.ID = accountIdGenerator.Uint()
	return err
}
