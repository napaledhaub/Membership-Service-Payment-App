package controllers

import (
	"net/http"
	"paymentapp/services"

	"github.com/gin-gonic/gin"
)

type BillingController struct {
	PaymentService *services.PaymentService
}

type BillRequest struct {
	ExpectedBillAmount float64 `json:"expected_bill_amount"`
}

func (ctrl *BillingController) VerifyBillAmount(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	var billRequest BillRequest
	if err := c.ShouldBindJSON(&billRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	if err := ctrl.PaymentService.VerifyBillAmount(authToken, billRequest.ExpectedBillAmount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bill verification failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill verification successful"})
}
