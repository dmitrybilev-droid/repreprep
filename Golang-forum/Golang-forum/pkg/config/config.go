package config

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	AuthGRPCAddr string
	ChatGRPCAddr string
	JWTSecret    string
}

func LoadUserConfig() *Config {

	return &Config{
		DBHost:       "localhost",
		DBPort:       "5432",
		DBUser:       "postgres",
		DBPassword:   "1234",
		DBName:       "forum_user_db",
		AuthGRPCAddr: "localhost:50051",
		ChatGRPCAddr: "localhost:50052",
		JWTSecret:    "12345token",
	}
}

func LoadChatConfig() *Config {

	return &Config{
		DBHost:       "localhost",
		DBPort:       "5432",
		DBUser:       "postgres",
		DBPassword:   "1234",
		DBName:       "forum_chat_db",
		AuthGRPCAddr: "localhost:50051",
		ChatGRPCAddr: "localhost:50052",
		JWTSecret:    "12345token",
	}
}
