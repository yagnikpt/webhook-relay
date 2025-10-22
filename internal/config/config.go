package config

import "os"

type Config struct {
	Port           string
	Environment    string
	GitHubClientID string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		Environment:    getEnv("ENVIRONMENT", "production"),
		GitHubClientID: getEnv("GITHUB_CLIENT_ID", "Ov23licjIRxaBMP2jm14"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
