package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIAPIKey string
	Port         string
	LogLevel     string
}

func LoadEnv(filename string) error {
	fmt.Printf("Attempting to load environment from file: %s\n", filename)

	err := godotenv.Load(filename)
	if err != nil {
		fmt.Printf("Warning: %s - %v\n", filename, err)
		return nil
	}

	fmt.Printf("Successfully loaded environment from file: %s\n", filename)
	return nil
}

func LoadConfig() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Warning: Could not get current working directory: %v\n", err)
	} else {
		fmt.Printf("Current working directory: %s\n", cwd)
	}

	possiblePaths := []string{
		".env",
		"../.env",
		"../../.env",
		os.Getenv("HOME") + "/TutorGPT/.env",
	}

	err = godotenv.Load(possiblePaths...)
	if err != nil {
		fmt.Println("Could not find .env file in any of the expected locations")
	} else {
		fmt.Println("Successfully loaded environment variables from .env file")
	}

	config := &Config{
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
		Port:         os.Getenv("PORT"),
		LogLevel:     os.Getenv("LOG_LEVEL"),
	}

	if config.Port == "" {
		config.Port = "8000"
	}

	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	if config.OpenAIAPIKey == "" {
		return nil, errors.New("OPENAI_API_KEY is required")
	}

	return config, nil
}

func PrintConfig(config *Config) {
	fmt.Println("Configuration:")
	fmt.Println("  Port:", config.Port)
	fmt.Println("  Log Level:", config.LogLevel)
}
