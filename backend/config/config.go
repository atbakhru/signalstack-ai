package config

import "os"

type Config struct {
	ServerAddress        string
	DatabaseURL          string
	JWTSecret            string
	AWSRegion            string
	S3Bucket             string
	GuardianAPIKey       string
	OpenAIAPIKey         string
	OpenAIEmbeddingModel string
	OpenAIChatModel      string
}

func Load() Config {
	return Config{
		ServerAddress:        getEnv("SERVER_ADDRESS", ":8080"),
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://signalstack:signalstack@localhost:5432/signalstack?sslmode=disable"),
		JWTSecret:            getEnv("JWT_SECRET", "change-me"),
		AWSRegion:            getEnv("AWS_REGION", "us-east-1"),
		S3Bucket:             getEnv("AWS_S3_BUCKET", "signalstack-raw-responses"),
		GuardianAPIKey:       getEnv("GUARDIAN_API_KEY", ""),
		OpenAIAPIKey:         getEnv("OPENAI_API_KEY", ""),
		OpenAIEmbeddingModel: getEnv("OPENAI_EMBEDDING_MODEL", "text-embedding-3-small"),
		OpenAIChatModel:      getEnv("OPENAI_CHAT_MODEL", "gpt-4o-mini"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
