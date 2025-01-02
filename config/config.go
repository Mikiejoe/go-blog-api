package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	MongoUrl   string
	DBName     string
	JWTSecret string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", ":8080"),
		MongoUrl:   getEnv("MONGO_URI", "localhost:27017"),
		DBName:     getEnv("DB_NAME", "blogger"),
		JWTSecret: getEnv("JWT_SECRET","dgfhjkluhjfjgjkfdnjkgjkdfhgf"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

type Collection struct{
	Users string
	Blogs string
	Likes string 
	Dislikes string
	Comments string
}

var Collections = initCollections()

func initCollections() Collection{
	return Collection{
		Users: "users",
		Blogs: "blogs",
		Likes: "likes",
		Dislikes: "dislikes",
		Comments: "comments",
	} 
}
