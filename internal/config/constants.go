package config

// this will be used to store all the hard coded value or magic number in code

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	JWTSecret   string
}

type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
}
