package controllers

import (
	"net/http"
	"paymentapp/models"
	"paymentapp/services"

	"github.com/gin-gonic/gin"
)

type ServiceMenuController struct {
	ServiceMenuService *services.ServiceMenuService
}

func (ctrl *ServiceMenuController) GetAllServiceMenus(c *gin.Context) {
	serviceMenus, err := ctrl.ServiceMenuService.GetAllServiceMenus()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get service menus: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, serviceMenus)
}

func (ctrl *ServiceMenuController) AddServiceMenu(c *gin.Context) {
	var serviceMenu models.ServiceMenu
	if err := c.ShouldBindJSON(&serviceMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input: " + err.Error()})
		return
	}
	if err := ctrl.ServiceMenuService.AddServiceMenu(&serviceMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add service menu: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, serviceMenu)
}
