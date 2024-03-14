package services

import (
	"eshell/dao/repo"
	"eshell/services/models"

	"github.com/universalmacro/common/singleton"
)

func GetDashboardService() *DashboardService {
	return dashboardServiceSingleton.Get()
}

var dashboardServiceSingleton = singleton.SingletonFactory(newDashboardService, singleton.Eager)

type DashboardService struct {
	dashboardRepo *repo.DashboardRepository
}

func newDashboardService() *DashboardService {
	return &DashboardService{dashboardRepo: repo.GetDashboardRepo()}
}

func (s *DashboardService) GetById(id uint) *models.Dashboard {
	dashboard, _ := s.dashboardRepo.GetById(id)
	if dashboard == nil {
		return nil
	}
	return &models.Dashboard{Dashboard: dashboard}
}

func (s *DashboardService) ListDashboard() []models.Dashboard {
	dashboards, _ := s.dashboardRepo.List()
	var result []models.Dashboard
	for i := range dashboards {
		result = append(result, models.Dashboard{Dashboard: &dashboards[i]})
	}
	return result
}
