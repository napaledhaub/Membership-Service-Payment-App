package models

type ServiceMenu struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	Name              string         `json:"name"`
	PricePerSession   float64        `json:"price_per_session"`
	TotalSessions     int            `json:"total_sessions"`
	Schedule          string         `json:"schedule"`
	DurationInMinutes int            `json:"duration_in_minutes"`
	Exercises         []Exercise     `json:"exercise_list" gorm:"foreignKey:ServiceMenuID"`
	Subscriptions     []Subscription `json:"subscription" gorm:"foreignKey:ServiceMenuID"`
}
