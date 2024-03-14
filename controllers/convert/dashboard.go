package convert

import (
	"encoding/json"
	apiModels "eshell/controllers/models"
	"eshell/dao/entities"
	models "eshell/services/models"

	"github.com/universalmacro/common/utils"
)

func ConvertDashboard(dashboard *models.Dashboard) *apiModels.Dashboard {
	return &apiModels.Dashboard{
		ID:   utils.UintToString(dashboard.ID()),
		Name: dashboard.Name,
	}
}

func ConvertVisualization(visualization apiModels.Visualization) entities.Visualization {
	var result entities.Visualization
	result.DSL, _ = json.Marshal(visualization.DSL)
	result.Type = visualization.Type
	result.Views, _ = json.Marshal(visualization.Views)
	return result
}

func ConvertBackVisualization(visualization entities.Visualization) apiModels.Visualization {
	var result apiModels.Visualization
	result.Type = visualization.Type
	json.Unmarshal(visualization.DSL, &result.DSL)
	json.Unmarshal(visualization.Views, &result.Views)
	return result
}
