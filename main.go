package main

import (
	"fmt"
	"log"
	"paymentapp/config"
	"paymentapp/controllers"
	"paymentapp/middleware"
	"paymentapp/models"
	"paymentapp/services"
	"paymentapp/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.LoggingMiddleware())

	initializeServicesAndControllers(r, db, cfg)

	r.Run(":8080")
}

func setupDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func initializeServicesAndControllers(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	db.AutoMigrate(
		&models.ServiceMenu{},
		&models.Exercise{},
		&models.Participant{},
		&models.AuthToken{},
		&models.Subscription{})

	emailService := &utils.EmailService{
		SMTPHost: cfg.Email.SMTPHost,
		SMTPPort: cfg.Email.SMTPPort,
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
	}

	serviceMenuService := &services.ServiceMenuService{DB: db}
	participantService := &services.ParticipantService{
		DB:            db,
		EncryptionKey: []byte(cfg.EncryptionKey),
	}
	passwordResetService := &services.PasswordResetService{
		DB:           db,
		EmailService: emailService,
	}
	registrationService := &services.RegistrationService{
		EmailService:  emailService,
		EncryptionKey: []byte(cfg.EncryptionKey),
	}
	subscriptionService := &services.SubscriptionService{DB: db}
	paymentService := &services.PaymentService{
		EmailService: emailService,
		DB:           db,
	}

	serviceMenuController := &controllers.ServiceMenuController{ServiceMenuService: serviceMenuService}
	authenticationService := &services.AuthenticationService{DB: db}
	authenticationController := &controllers.AuthenticationController{
		AuthenticationService: authenticationService,
		ParticipantService:    participantService,
		PasswordResetService:  passwordResetService,
		RegistrationService:   registrationService,
	}
	infoManagementController := &controllers.InfoManagementController{
		ParticipantService: participantService,
	}
	subscriptionController := &controllers.SubscriptionController{
		SubscriptionService: subscriptionService,
	}
	billingController := &controllers.BillingController{
		PaymentService: paymentService,
	}
	participantController := &controllers.ParticipantController{
		PaymentService: paymentService,
	}

	r.GET("/service-menu/list", serviceMenuController.GetAllServiceMenus)
	r.POST("/service-menu/add", serviceMenuController.AddServiceMenu)

	r.POST("/auth/register", authenticationController.Register)
	r.POST("/auth/check-status", authenticationController.CheckStatus)
	r.POST("/auth/login", authenticationController.Login)
	r.POST("/auth/confirm", authenticationController.Confirm)
	r.POST("/auth/refresh", authenticationController.Refresh)
	r.POST("/auth/forgot-password", authenticationController.ForgotPassword)
	r.POST("/auth/forgot-password-confirmation", authenticationController.ForgotPasswordConfirmation)
	r.POST("/auth/logout", authenticationController.Logout)

	r.POST("/info-management/update-fullname", infoManagementController.UpdateFullName)
	r.POST("/info-management/update-credit-card-info", infoManagementController.UpdateCreditCardInfo)
	r.POST("/info-management/update-password", infoManagementController.UpdatePassword)

	r.GET("subscription/subscribe", subscriptionController.SubscribeToService)
	r.GET("subscription/list", subscriptionController.GetSubscriptionList)
	r.GET("subscription/extend", subscriptionController.ExtendSubscription)
	r.GET("subscription/cancel", subscriptionController.CancelSubscription)

	r.POST("/billing/verify-bill", billingController.VerifyBillAmount)

	r.POST("/participant/send-email-otp", participantController.SendEmailVerification)
	r.POST("/participant/is-payment-otp-expired", participantController.IsPaymentOTPExpired)
	r.POST("/participant/verify-payment", participantController.VerifyPayment)
}
