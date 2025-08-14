package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Kafka    KafkaConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type ServerConfig struct {
	Host string
	Port string
}

// Загрузка конфигурации из .env файла
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{os.Getenv("KAFKA_BROKERS")},
			Topic:   os.Getenv("KAFKA_TOPIC"),
			GroupID: os.Getenv("KAFKA_GROUP_ID"),
		},
		Server: ServerConfig{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
	}
	return config, nil
}

// GetDatabaseUrl формирует строку подключения к базе данных
func (c *Config) GetDatabaseUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode)
}

// GetServerAddress формирует адрес сервера
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
