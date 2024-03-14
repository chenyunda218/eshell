package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/universalmacro/common/snowflake"
	"gorm.io/gorm"
)

type Dashboard struct {
	gorm.Model
	AccountID      uint `gorm:"index"`
	Name           string
	Visualizations VisualizationList
}

type VisualizationList []Visualization

func (v *VisualizationList) Scan(value any) error {
	return json.Unmarshal(value.([]byte), v)
}

func (v VisualizationList) Value() (driver.Value, error) {
	return json.Marshal(v)
}

type Visualization struct {
	Type  string `json:"type"`
	DSL   JSON   `json:"dsl"`
	Views JSON   `json:"views"`
}

var dashboardIdGenerator = snowflake.NewIdGenertor(0)

func (a *Dashboard) BeforeCreate(tx *gorm.DB) error {
	a.Model.ID = dashboardIdGenerator.Uint()
	return nil
}

type JSON json.RawMessage

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
