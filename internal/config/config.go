package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	APIID          int32
	APIHash        string
	BotToken       string
	MongoURI       string
	OwnerID        int64
	SudoUsers      []int64
	LoggerID       int64
	StringSessions []string
	BotUsername    string
	DurationLimit  int
	QueueLimit     int
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("[config] required env var %s is missing — check your .env file", key)
	}
	return v
}

func envInt64(key string, fallback int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatalf("[config] %s must be an integer, got %q", key, v)
	}
	return n
}

func envInt(key string, fallback int) int {
	return int(envInt64(key, int64(fallback)))
}

func envIDList(key string) []int64 {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			log.Fatalf("[config] %s contains an invalid id %q", key, p)
		}
		out = append(out, n)
	}
	return out
}

func envStrList(key string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func Load() *Config {
	apiID, err := strconv.Atoi(mustEnv("API_ID"))
	if err != nil {
		log.Fatalf("[config] API_ID must be numeric: %v", err)
	}

	cfg := &Config{
		APIID:          int32(apiID),
		APIHash:        mustEnv("API_HASH"),
		BotToken:       mustEnv("TOKEN"),
		MongoURI:       mustEnv("MONGO_DB_URI"),
		OwnerID:        envInt64("OWNER_ID", 0),
		SudoUsers:      envIDList("SUDO_USERS"),
		LoggerID:       envInt64("LOGGER_ID", 0),
		StringSessions: envStrList("STRING_SESSIONS"),
		BotUsername:    strings.TrimPrefix(os.Getenv("BOT_USERNAME"), "@"),
		DurationLimit:  envInt("DURATION_LIMIT", 0),
		QueueLimit:     envInt("QUEUE_LIMIT", 50),
	}

	if cfg.BotUsername == "" {
		cfg.BotUsername = "Kokomimusicbot"
	}

	return cfg
}
