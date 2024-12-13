package controllers

import (
	"net/http"
	"paymentapp/services"

	"github.com/gin-gonic/gin"
)

type ParticipantController struct {
	PaymentService *services.PaymentService
}

func (ctrl *ParticipantController) SendEmailVerification(c *gin.Context) {
	var emailRequest EmailRequest
	participantId := c.Query("participantId")

	if err := c.ShouldBindJSON(&emailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	err := ctrl.PaymentService.SendEmailVerification(emailRequest.Email, participantId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to send verification email: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email has been sent"})
}

func (ctrl *ParticipantController) VerifyPayment(c *gin.Context) {
	var verificationRequest VerificationRequest

	if err := c.ShouldBindJSON(&verificationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	if err := ctrl.PaymentService.VerifyPaymentOTP(verificationRequest.Email, verificationRequest.OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment verification failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment verification successful"})
}

func (ctrl *ParticipantController) IsPaymentOTPExpired(c *gin.Context) {
	var emailRequest EmailRequest

	if err := c.ShouldBindJSON(&emailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	expired := ctrl.PaymentService.IsPaymentOTPExpired(emailRequest.Email)
	if !expired {
		c.JSON(http.StatusOK, gin.H{"message": "Payment OTP is still valid"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Payment OTP is not valid"})
	}
}
