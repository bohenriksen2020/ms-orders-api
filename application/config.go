package application

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RedisAddress string
	PostgresDSN  string
	DatabaseType string // "redis" or "postgres"
	ServerPort   uint16
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	cfg := Config{
		RedisAddress: "localhost:6379",                                        // Default Redis address
		PostgresDSN:  "postgres://user:pass@localhost/dbname?sslmode=disable", // Default Postgres DSN
		DatabaseType: "redis",                                                 // Default database is Redis
		ServerPort:   3000,                                                    // Default server port
	}

	// Load Redis address from environment
	if redisAddr, exists := os.LookupEnv("REDIS_ADDRESS"); exists {
		cfg.RedisAddress = redisAddr
	}

	// Load Postgres DSN from environment
	if postgresDSN, exists := os.LookupEnv("POSTGRES_DSN"); exists {
		cfg.PostgresDSN = postgresDSN
	}

	// Load database type from environment (either "redis" or "postgres")
	if dbType, exists := os.LookupEnv("DATABASE_TYPE"); exists {
		cfg.DatabaseType = dbType
	}

	// Load server port from environment
	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
			cfg.ServerPort = uint16(port)
		}
	}

	fmt.Println("ServerPort: ", cfg.ServerPort)
	fmt.Println("RedisAddress: ", cfg.RedisAddress)
	fmt.Println("PostgresDSN: ", cfg.PostgresDSN)
	fmt.Println("DatabaseType: ", cfg.DatabaseType)

	return cfg
}
