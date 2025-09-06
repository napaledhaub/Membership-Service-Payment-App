package config

import (
	"encoding/json"
	"os"
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
		SMTPHost     string `json:"smtp_host"`
		SMTPPort     int    `json:"smtp_port"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		IsSMTPActive bool   `json:"is_smpt_active"`
	} `json:"email"`
	EncryptionKey string `json:"encryption_key"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
