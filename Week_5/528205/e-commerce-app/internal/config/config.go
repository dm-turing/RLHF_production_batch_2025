package config

type Config struct {
	Port        string
	DatabaseDSN string
}

func NewConfig() *Config {
	return &Config{
		Port:        "8080",
		DatabaseDSN: "file:user_management.db?cache=shared&mode=rwc", // Use an in-memory SQLite database for testing purposes.
	}
}
