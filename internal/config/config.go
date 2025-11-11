// Package config handles all configurations
package config

import (
	"log"
	"os"
	"strconv"
)

const (
	ServerAddressKey = "SERVER_ADDRESS"
	ServerPortKey    = "SERVER_PORT"
)

// Config holds all configuration parameters
type Config struct {
	ServerAddress string
	ServerPort    int
}

// InitConfig initializes the configurations parameters from all sources
func InitConfig() Config {
	serverAddress := os.Getenv(ServerAddressKey)

	port, err := strconv.Atoi(os.Getenv(ServerPortKey))
	if err != nil {
		log.Panicf("[ENV] Invalid server port: %v", err)
	}

	return Config{
		ServerAddress: serverAddress,
		ServerPort:    port,
	}
}
