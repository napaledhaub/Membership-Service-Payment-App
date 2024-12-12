package models

type Exercise struct {
	ID                uint        `json:"id" gorm:"primaryKey"`
	Name              string      `json:"name"`
	DurationInMinutes int         `json:"duration_in_minutes"`
	Description       string      `json:"description"`
	ServiceMenuID     uint        `json:"service_menu_id"`
	ServiceMenu       ServiceMenu `gorm:"foreignKey:ServiceMenuID" json:"service_menu,omitempty"`
}
