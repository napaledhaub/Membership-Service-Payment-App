package controllers

import (
	"net/http"
	"paymentapp/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	SubscriptionService *services.SubscriptionService
}

func (ctrl *SubscriptionController) GetSubscriptionList(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	subscriptions, err := ctrl.SubscriptionService.GetSubscriptionList(authToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Get subscription list failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

func (ctrl *SubscriptionController) SubscribeToService(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	serviceId := c.Query("serviceId")
	serviceMenuID, err := strconv.ParseUint(serviceId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	subscription, err := ctrl.SubscriptionService.SubscribeToService(authToken, uint(serviceMenuID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Subscribe to this service failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription successful", "subscription": subscription})
}

func (ctrl *SubscriptionController) CancelSubscription(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	subscriptionId := c.Query("subscriptionId")
	subscriptionID, err := strconv.ParseUint(subscriptionId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	if err = ctrl.SubscriptionService.CancelSubscription(authToken, uint(subscriptionID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cancel subscription failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
}

func (ctrl *SubscriptionController) ExtendSubscription(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	subscriptionId := c.Query("subscriptionId")
	subscriptionID, err := strconv.ParseUint(subscriptionId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	if err = ctrl.SubscriptionService.ExtendSubscription(authToken, uint(subscriptionID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Session extension failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session duration successfully extended"})
}
