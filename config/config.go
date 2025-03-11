package config

import (
	"fmt"
	"log"

	"myapp/models"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.AutomaticEnv()

	// Check if running in Docker by checking environment variables
	dbHost := viper.GetString("DB_HOST")
	redisHost := viper.GetString("REDIS_HOST")
	rabbitmqHost := viper.GetString("RABBITMQ_HOST")

	// Server configuration
	if dbHost != "" {
		// In Docker, use port 3000
		viper.SetDefault("server.port", 3000)
	} else {
		// Locally, use port 8000
		viper.SetDefault("server.port", 8000)
	}
	viper.SetDefault("server.host", "0.0.0.0")

	// Database configuration
	if dbHost != "" {
		viper.Set("database.host", dbHost)
	} else {
		viper.Set("database.host", "localhost")
	}

	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "myapp")

	// Bind database environment variables
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")

	// Redis configuration
	if redisHost != "" {
		viper.Set("redis.host", redisHost)
	} else {
		viper.Set("redis.host", "localhost")
	}
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Bind Redis environment variables
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")

	// RabbitMQ configuration
	if rabbitmqHost != "" {
		viper.Set("rabbitmq.host", rabbitmqHost)
	} else {
		viper.Set("rabbitmq.host", "localhost")
	}
	viper.SetDefault("rabbitmq.port", "5672")
	viper.SetDefault("rabbitmq.user", "guest")
	viper.SetDefault("rabbitmq.password", "guest")

	// Bind RabbitMQ environment variables
	viper.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
	viper.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
	viper.BindEnv("rabbitmq.user", "RABBITMQ_USER")
	viper.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")

	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiresIn", "24h")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func ConnectDB(config *DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		"disable",
		"UTC",
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&models.User{}, &models.Product{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Connected Successfully to Database")
	return db
}
