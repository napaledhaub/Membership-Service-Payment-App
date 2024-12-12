package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"paymentapp/controllers"
	"paymentapp/middleware"
	"paymentapp/models"
	"paymentapp/services"
	"paymentapp/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Database struct {
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		Port     int    `json:"port"`
		SSLMode  string `json:"sslmode"`
	} `json:"database"`
	Email struct {
		SMTPHost string `json:"smtp_host"`
		SMTPPort int    `json:"smtp_port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"email"`
	EncryptionKey string `json:"encryption_key"`
}

func loadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(
		&models.ServiceMenu{},
		&models.Exercise{},
		&models.Participant{},
		&models.AuthToken{},
		&models.Subscription{})

	r := gin.Default()

	r.Use(middleware.LoggingMiddleware())

	emailService := &utils.EmailService{
		SMTPHost: cfg.Email.SMTPHost,
		SMTPPort: cfg.Email.SMTPPort,
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
	}

	serviceMenuService := &services.ServiceMenuService{DB: db}
	serviceMenuController := &controllers.ServiceMenuController{ServiceMenuService: serviceMenuService}

	participantService := &services.ParticipantService{DB: db}
	passwordResetService := &services.PasswordResetService{
		DB:           db,
		EmailService: emailService,
	}
	authenticationService := &services.AuthenticationService{
		DB:            db,
		EmailService:  emailService,
		EncryptionKey: []byte("1234567812345678"),
	}
	authenticationController := &controllers.AuthenticationController{
		AuthenticationService: authenticationService,
		ParticipantService:    participantService,
		PasswordResetService:  passwordResetService,
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

	r.Run(":8080")
}
