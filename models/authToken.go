package models

import "time"

type AuthToken struct {
	ID                 uint        `gorm:"primaryKey" json:"id"`
	Token              string      `json:"token"`
	ExpirationDateTime time.Time   `json:"expiration_date_time"`
	ParticipantID      uint        `json:"participant_id"`
	Participant        Participant `gorm:"foreignKey:ParticipantID" json:"participant"`
}

func (token *AuthToken) UpdateExpiration() {
	token.ExpirationDateTime = time.Now().Add(24 * time.Hour)
}
