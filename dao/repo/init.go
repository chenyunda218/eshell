package repo

import (
	"eshell/dao/entities"
	"eshell/ioc"
)

func init() {
	ioc.GetDBInstance().AutoMigrate(&entities.Account{}, &entities.Dashboard{})
}
