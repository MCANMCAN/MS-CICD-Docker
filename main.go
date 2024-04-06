package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
type User struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var (
	db          *sql.DB
	redisClient *redis.Client
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize PostgreSQL database connection
	dbConfig := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     atoi(os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}
	connectToDB(&dbConfig)

	// Initialize Redis client
	redisConfig := RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // Default DB
	}
	connectToRedis(&redisConfig)

	// Define HTTP routes
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/user/", getUserHandler)

	// Start the HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %v", err)
	}
	return i
}

func connectToDB(config *DBConfig) {
	// Create the PostgreSQL connection string
	connStr := "host=" + config.Host + " port=" + strconv.Itoa(config.Port) +
		" user=" + config.User + " password=" + config.Password +
		" dbname=" + config.DBName + " sslmode=" + config.SSLMode

	// Open a database connection
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Check if the connection is established
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}
	log.Println("Connected to the database")
}

func connectToRedis(config *RedisConfig) {
	// Create a new Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	// Ping the Redis server to check the connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis: ", err)
	}
	log.Println("Connected to Redis")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Insert the new user into the database
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
		user.Username, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		log.Println("Error registering user:", err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Implement login logic here
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Implement get user logic here
}
