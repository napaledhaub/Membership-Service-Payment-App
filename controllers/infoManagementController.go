package controllers

import (
	"net/http"
	"paymentapp/models"
	"paymentapp/services"

	"github.com/gin-gonic/gin"
)

type UpdateInfoManagementRequest struct {
	NewFullName       string            `json:"new_full_name"`
	NewCreditCardInfo models.CreditCard `json:"new_credit_card_info"`
	OldPassword       string            `json:"old_password"`
	NewPassword       string            `json:"new_password"`
}

type InfoManagementController struct {
	ParticipantService *services.ParticipantService
}

func (ctrl *InfoManagementController) UpdateFullName(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	var updateInfoRequest UpdateInfoManagementRequest
	if err := c.ShouldBindJSON(&updateInfoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	err := ctrl.ParticipantService.UpdateFullName(authToken, updateInfoRequest.NewFullName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Update full name failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Full name successfully updated"})
}

func (ctrl *InfoManagementController) UpdateCreditCardInfo(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	var request UpdateInfoManagementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	err := ctrl.ParticipantService.UpdateCreditCardInfo(authToken, request.NewCreditCardInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Update credit card failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credit card information successfully updated"})
}

func (ctrl *InfoManagementController) UpdatePassword(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	var request UpdateInfoManagementRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	err := ctrl.ParticipantService.UpdatePassword(authToken, request.OldPassword, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Update password failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully updated"})
}
