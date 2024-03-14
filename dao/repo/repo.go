package repo

import (
	"eshell/dao/entities"
	"eshell/ioc"

	"github.com/universalmacro/common/dao"
	"github.com/universalmacro/common/singleton"
)

func GetAccountRepo() *AccountRepository {
	return accountRepo.Get()
}

var accountRepo = singleton.SingletonFactory(createAccountRepository, singleton.Eager)

type AccountRepository struct {
	*dao.Repository[entities.Account]
}

func createAccountRepository() *AccountRepository {
	return &AccountRepository{
		Repository: &dao.Repository[entities.Account]{DB: ioc.GetDBInstance()},
	}
}

func GetDashboardRepo() *DashboardRepository {
	return dashboardRepo.Get()
}

var dashboardRepo = singleton.SingletonFactory(createDashboardRepository, singleton.Eager)

type DashboardRepository struct {
	*dao.Repository[entities.Dashboard]
}

func createDashboardRepository() *DashboardRepository {
	return &DashboardRepository{
		Repository: &dao.Repository[entities.Dashboard]{DB: ioc.GetDBInstance()},
	}
}
