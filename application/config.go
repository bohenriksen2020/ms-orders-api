package application

import (
	"strconv"
	"os"
	"fmt"
)

type Config struct {
	RedisAddress string
	ServerPort  uint16
}

func LoadConfig() Config {
	cfg := Config {
		RedisAddress: "localhost:6379",
		ServerPort: 3000,

	}
	if redisAddr, exists := os.LookupEnv("REDIS_ADDRESS"); exists {
		cfg.RedisAddress = redisAddr
	}
	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
			
		}
		
	}

	fmt.Println("ServerPort: ", cfg.ServerPort)
	fmt.Println("RedisAddress: ", cfg.RedisAddress)
	return cfg
}