package models

import (
	"time"
)

type Subscription struct {
	ID                uint        `gorm:"primaryKey" json:"id"`
	StartDate         time.Time   `json:"start_date"`
	EndDate           time.Time   `json:"end_date"`
	RemainingSessions int         `json:"remaining_sessions"`
	ParticipantID     uint        `json:"participant_id"`
	Participant       Participant `gorm:"foreignKey:ParticipantID" json:"participant,omitempty"`
	ServiceMenuID     uint        `json:"service_menu_id"`
	ServiceMenu       ServiceMenu `gorm:"foreignKey:ServiceMenuID" json:"service_menu,omitempty"`
}
