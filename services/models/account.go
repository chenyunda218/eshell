package models

import (
	"eshell/dao/entities"
	"eshell/dao/repo"
)

type Account struct {
	*entities.Account
}

func (a *Account) CreateDashbaord(name string) *Dashboard {
	dashboard := &entities.Dashboard{Name: name, AccountID: a.ID}
	repo.GetDashboardRepo().Create(dashboard)
	return &Dashboard{Dashboard: dashboard}
}

type Dashboard struct {
	*entities.Dashboard
}

func (d *Dashboard) ID() uint {
	return d.Dashboard.ID
}

func (d *Dashboard) UpdateVisualizations(visualizations ...entities.Visualization) *Dashboard {
	d.Dashboard.Visualizations = visualizations
	return d
}

func (d *Dashboard) Submit() *Dashboard {
	repo.GetDashboardRepo().Update(d.Dashboard)
	return d
}

func (d *Dashboard) Delete() *Dashboard {
	repo.GetDashboardRepo().Delete(d.Dashboard)
	return d
}
