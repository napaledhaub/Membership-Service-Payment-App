package controllers

import (
	"net/http"
	"paymentapp/models"
	"paymentapp/services"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	AuthenticationService *services.AuthenticationService
	ParticipantService    *services.ParticipantService
	PasswordResetService  *services.PasswordResetService
	RegistrationService   *services.RegistrationService
}

type EmailRequest struct {
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerificationRequest struct {
	EmailRequest
	OTP string `json:"otp"`
}

type ForgotPasswordRequest struct {
	ResetToken  string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}

func (ctrl *AuthenticationController) Register(c *gin.Context) {
	var registerRequest models.Participant
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	if err := ctrl.RegistrationService.Register(registerRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func (ctrl *AuthenticationController) CheckStatus(c *gin.Context) {
	var emailRequest EmailRequest
	if err := c.ShouldBindJSON(&emailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	participant, err := ctrl.ParticipantService.FindByEmail(emailRequest.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Check status failed: " + err.Error()})
		return
	}

	if participant == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Membership status: " + string(models.NOT_REGISTERED)})
		return
	}
	if participant.IsVerified {
		c.JSON(http.StatusOK, gin.H{"message": "Membership status: " + string(models.REGISTERED)})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Membership status: " + string(models.NOT_VALIDATED)})
	}
}

func (ctrl *AuthenticationController) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	token, err := ctrl.AuthenticationService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Login failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": string(token)})
}

func (ctrl *AuthenticationController) Logout(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	if err := ctrl.AuthenticationService.Logout(authToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Logout failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ctrl *AuthenticationController) Confirm(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	var verificationRequest VerificationRequest
	if err := c.ShouldBindJSON(&verificationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	if err := ctrl.AuthenticationService.ValidateOtp(authToken, verificationRequest.OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Confirm email failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your account has been successfully validated. Membership status: " + string(models.REGISTERED)})
}

func (ctrl *AuthenticationController) Refresh(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	newAuthToken, err := ctrl.AuthenticationService.RefreshToken(authToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Refresh failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, newAuthToken)
}

func (ctrl *AuthenticationController) ForgotPassword(c *gin.Context) {
	var emailRequest EmailRequest
	if err := c.ShouldBindJSON(&emailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	if err := ctrl.PasswordResetService.RequestPasswordReset(emailRequest.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password reset request failed to send: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email has been successfully sent. Please check your email"})
}

func (ctrl *AuthenticationController) ForgotPasswordConfirmation(c *gin.Context) {
	var forgotPasswordRequest ForgotPasswordRequest
	if err := c.ShouldBindJSON(&forgotPasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request: " + err.Error()})
		return
	}

	if err := ctrl.PasswordResetService.ConfirmPasswordReset(forgotPasswordRequest.ResetToken, forgotPasswordRequest.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password change failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully changed"})
}
