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

type Config struct {
	Env           string
	ServerAddress string
	ServerPort    int
}

func NewConfig() Config {
	serverAddress := os.Getenv(ServerAddressKey)

	port, err := strconv.Atoi(os.Getenv(ServerPortKey))
	if err != nil {
		log.Fatalf("[ENV] Invalid server port: %v", err)
	}

	return Config{
		ServerAddress: serverAddress,
		ServerPort:    port,
	}
}
