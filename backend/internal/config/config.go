package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port             string
	AWSRegion        string
	AWSEndpoint      string
	DynamoUserTable  string
	DynamoOrderTable string
	UserServiceURL   string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:             getEnv("PORT", "8080"),
		AWSRegion:        getEnv("AWS_REGION", "ap-northeast-2"),
		AWSEndpoint:      getEnv("AWS_ENDPOINT", ""),
		DynamoUserTable:  getEnv("DYNAMO_USER_TABLE", ""),
		DynamoOrderTable: getEnv("DYNAMO_ORDER_TABLE", ""),
		UserServiceURL:   getEnv("USER_SERVICE_URL", "http://localhost:8081"),
	}
	if cfg.DynamoUserTable == "" || cfg.DynamoOrderTable == "" {
		return nil, fmt.Errorf("DynamoDB 테이블 이름이 비어 있음")
	}
	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
