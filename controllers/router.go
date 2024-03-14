package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/server"
)

func Init() {
	r := gin.Default()
	accountController := NewAccountController()
	r.Use(server.CorsMiddleware())
	r.Use(auth())
	r.POST("/accounts", accountController.CreateAccount)
	r.POST("/sessions", accountController.CreateSession)
	r.POST("/search", accountController.Search)
	r.POST("/dashboards", accountController.CreateDashboard)
	r.GET("/dashboards", accountController.ListDashboard)
	r.DELETE("/dashboards/:dashboardId", accountController.DeleteDashboard)
	r.PUT("/dashboards/:dashboardId", accountController.UpdateDashboard)
	r.GET("/dashboards/:dashboardId", accountController.GetDashboard)
	r.PUT("/dashboards/:dashboardId/visualizations", accountController.UpdateDashboardVisualizations)
	r.GET("/dashboards/:dashboardId/visualizations", accountController.GetVisualizations)
	r.Run()
}
