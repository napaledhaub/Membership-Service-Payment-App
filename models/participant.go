package models

import (
	"time"
)

type Participant struct {
	ID                   uint           `json:"id" gorm:"primaryKey"`
	Name                 string         `json:"name"`
	Email                string         `gorm:"unique" json:"email"`
	Password             string         `json:"password"`
	PhoneNumber          string         `json:"phone_number"`
	IsVerified           bool           `json:"is_verified"`
	VerificationCode     string         `json:"verification_code"`
	ResetToken           string         `json:"reset_token"`
	PaymentOTP           string         `json:"payment_otp"`
	PaymentOTPExpiration time.Time      `json:"payment_otp_expiration"`
	BillAmount           float64        `json:"bill_amount"`
	IsBillVerified       bool           `json:"is_bill_verified"`
	CreditCardInfo       string         `json:"credit_card_info"`
	Subscriptions        []Subscription `json:"subscriptions" gorm:"foreignKey:ParticipantID"`
}

type ParticipantStatus string

const (
	REGISTERED     ParticipantStatus = "REGISTERED"
	NOT_VALIDATED  ParticipantStatus = "NOT VALIDATED"
	NOT_REGISTERED ParticipantStatus = "NOT REGISTERED"
)

type CreditCard struct {
	CardNo      string `json:"card_no"`
	Cvv         string `json:"cvv"`
	ExpiredDate string `json:"expired_date"`
	OwnerName   string `json:"owner_name"`
}
