package controllers

import (
	"eshell/client"
	"eshell/controllers/convert"
	apiModels "eshell/controllers/models"
	"eshell/dao/entities"
	"eshell/ioc"
	"eshell/ioc/config"
	"eshell/services"
	"eshell/services/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/universalmacro/common/utils"
)

type AccountController struct {
	accountServices  *services.AccountService
	dashboardService *services.DashboardService
	esClient         *client.Client
}

func NewAccountController() *AccountController {
	return &AccountController{
		accountServices:  services.GetAccountService(),
		dashboardService: services.GetDashboardService(),
		esClient:         ioc.GetEsClient()}
}

type Headers struct {
	Authorization string `header:"Authorization"`
	ApiKey        string `header:"ApiKey"`
}

func (ac *AccountController) CreateAccount(ctx *gin.Context) {
	var headers Headers
	ctx.ShouldBindHeader(&headers)
	apiKey := config.GetConfig().GetString("apiKey")
	if headers.ApiKey != apiKey {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var createAccountRequest apiModels.CreateAccountRequest
	ctx.ShouldBindJSON(&createAccountRequest)
	account := ac.accountServices.Create(createAccountRequest.Account, createAccountRequest.Password)
	if account == nil {
		ctx.JSON(400, gin.H{"error": "account already exists"})
		return
	}
	ctx.JSON(201, nil)
}

func (ac *AccountController) CreateSession(ctx *gin.Context) {
	var createSessionRequest apiModels.CreateSessionRequest
	ctx.ShouldBindJSON(&createSessionRequest)
	token, err := ac.accountServices.CreateSession(createSessionRequest.Account, createSessionRequest.Password)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "account not found"})
		return
	}
	ctx.JSON(200, gin.H{"token": token})
}

func (ac *AccountController) Search(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	index := ctx.Query("index")
	ac.esClient.Query(ctx, index)
}

func (ac *AccountController) CreateDashboard(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var createDashboardRequest apiModels.CreateDashboardRequest
	ctx.ShouldBindJSON(&createDashboardRequest)
	dashboard := account.CreateDashbaord(createDashboardRequest.Name)
	ctx.JSON(201, convert.ConvertDashboard(dashboard))
}

func (ac *AccountController) ListDashboard(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	dashboards := ac.dashboardService.ListDashboard()
	var apiDashboards []*apiModels.Dashboard = []*apiModels.Dashboard{}
	for i := range dashboards {
		apiDashboards = append(apiDashboards, convert.ConvertDashboard(&dashboards[i]))
	}
	ctx.JSON(200, apiDashboards)
}

func (ac *AccountController) UpdateDashboard(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	var updateDashboardRequest apiModels.UpdateDashboardRequest
	ctx.ShouldBindJSON(&updateDashboardRequest)
}

func (ac *AccountController) UpdateDashboardVisualizations(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	dashboardId := ctx.Param("dashboardId")
	dashboard := ac.dashboardService.GetById(utils.StringToUint(dashboardId))
	var visualizations []apiModels.Visualization
	ctx.ShouldBindJSON(&visualizations)
	var vs []entities.Visualization = []entities.Visualization{}
	for i := range visualizations {
		vs = append(vs, convert.ConvertVisualization(visualizations[i]))
	}
	dashboard.UpdateVisualizations(vs...).Submit()
	ctx.JSON(200, visualizations)
}

func (ac *AccountController) GetDashboard(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	dashboardId := ctx.Param("dashboardId")
	dashboard := ac.dashboardService.GetById(utils.StringToUint(dashboardId))
	ctx.JSON(200, dashboard)
}

func (ac *AccountController) GetVisualizations(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	dashboardId := ctx.Param("dashboardId")
	dashboard := ac.dashboardService.GetById(utils.StringToUint(dashboardId))
	if dashboard == nil {
		ctx.JSON(404, gin.H{"error": "dashboard not found"})
		return
	}
	var visualizations []apiModels.Visualization = []apiModels.Visualization{}
	for i := range dashboard.Visualizations {
		visualizations = append(visualizations, convert.ConvertBackVisualization(dashboard.Visualizations[i]))
	}
	ctx.JSON(200, visualizations)
}

func (ac *AccountController) DeleteDashboard(ctx *gin.Context) {
	account := getAccount(ctx)
	if account == nil {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	dashboardId := ctx.Param("dashboardId")
	dashboard := ac.dashboardService.GetById(utils.StringToUint(dashboardId))
	if dashboard == nil {
		ctx.JSON(404, gin.H{"error": "dashboard not found"})
		return
	}
	dashboard.Delete()
	ctx.JSON(204, nil)
}

func auth() func(ctx *gin.Context) {
	var accountService = services.GetAccountService()
	return func(ctx *gin.Context) {
		var headers Headers
		ctx.ShouldBindHeader(&headers)
		authorization := headers.Authorization
		splited := strings.Split(authorization, " ")
		if authorization != "" && len(splited) == 2 {
			account, err := accountService.VerifyToken(ctx, splited[1])
			if account != nil && err == nil {
				ctx.Set("account", account)
			} else {
				// TODO
			}
		} else {
			// TODO
		}
		ctx.Next()
	}
}

func getAccount(ctx *gin.Context) *models.Account {
	accountInterface, ok := ctx.Get("account")
	if !ok {
		return nil
	}
	account, ok := accountInterface.(*models.Account)
	if !ok {
		return nil
	}
	return account
}
