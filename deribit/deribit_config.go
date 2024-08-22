package deribit

import (
	"os"
)

type DeribitConfig struct {
	API_URL_BASE             string
	WS_URL_BASE              string
	WS_HOST                  string
	WS_PATH                  string
	GET_INSTRUMENTS_ENDPOINT string
	ORDERBOOK_INTERVAL       string
	KAFKA_SERVER_ADDRESS     string
	KAFKA_PRODUCER_PORT      string
	KAFKA_CONSUMER_PORT      string
}

func (cfg *DeribitConfig) readFromEnv() error {
	cfg.API_URL_BASE = getEnvOrDefault("API_URL_BASE", "https://test.deribit.com/api/v2")
	cfg.WS_URL_BASE = getEnvOrDefault("WS_URL_BASE", "wss://test.deribit.com/ws/api/v2")
	cfg.WS_HOST = getEnvOrDefault("WS_HOST", "test.deribit.com")
	cfg.WS_PATH = getEnvOrDefault("WS_PATH", "/ws/api/v2")
	cfg.GET_INSTRUMENTS_ENDPOINT = getEnvOrDefault("GET_INSTRUMENTS", "/public/get_instruments")
	cfg.ORDERBOOK_INTERVAL = getEnvOrDefault("ORDERBOOK_INTERVAL", "100ms")
	cfg.KAFKA_SERVER_ADDRESS = getEnvOrDefault("KAFKA_SERVER_ADDRESS", "localhost:9092")
	cfg.KAFKA_PRODUCER_PORT = getEnvOrDefault("KAFKA_PRODUCER_PORT", ":8080")
	cfg.KAFKA_CONSUMER_PORT = getEnvOrDefault("KAFKA_CONSUMER_PORT", ":8081")
	return nil
}

func getEnvOrDefault(envVarName, defaultVal string) string {
	ret, ok := os.LookupEnv(envVarName)
	if !ok {
		return defaultVal
	}
	return ret
}
