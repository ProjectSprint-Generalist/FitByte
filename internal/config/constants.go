package config

import "gorm.io/gorm"

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	JWTSecret   string
	DB          *gorm.DB
}

type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}
