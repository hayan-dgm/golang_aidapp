// package db

// import (
// 	"database/sql"
// 	"log"
// 	"os"

// 	_ "github.com/mattn/go-sqlite3"
// )

// var DB *sql.DB

// func InitDB() {
// 	var err error
// 	dbPath := os.Getenv("DATABASE_URL")
// 	if dbPath == "" {
// 		dbPath = "aid_app.db"
// 	}

// 	DB, err = sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	createTables()
// }

// func createTables() {
// 	tables := []string{
// 		`CREATE TABLE IF NOT EXISTS users (
// 			id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			username TEXT UNIQUE NOT NULL,
// 			password TEXT NOT NULL,
// 			isAdmin BOOLEAN DEFAULT FALSE
// 		)`,
// 		`CREATE TABLE IF NOT EXISTS families (...)`, // Full schema from earlier
// 		`CREATE TABLE IF NOT EXISTS active_sessions (...)`,
// 		`CREATE TABLE IF NOT EXISTS revoked_tokens (...)`,
// 		`CREATE TABLE IF NOT EXISTS logs (...)`,
// 	}

// 	for _, table := range tables {
// 		if _, err := DB.Exec(table); err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

// func CloseDB() {
// 	DB.Close()
// }

// ===============================================
// package db

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	_ "github.com/mattn/go-sqlite3"
// )

// // DB instance
// var DB *sql.DB

// // InitDB initializes the SQLite database and tables
// func InitDB() {
// 	var err error
// 	DB, err = sql.Open("sqlite3", "aid_app.db")
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	// Create tables
// 	createTables := []string{
// 		`CREATE TABLE IF NOT EXISTS users (
// 			id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			username TEXT NOT NULL,
// 			password TEXT NOT NULL,
// 			isAdmin BOOLEAN NOT NULL
// 		)`,
// 		`CREATE TABLE IF NOT EXISTS families (
// 			id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			fullName TEXT NOT NULL,
// 			nationalID TEXT NOT NULL,
// 			familyBookID TEXT NOT NULL,
// 			phoneNumber TEXT NOT NULL,
// 			familyMembers INTEGER NOT NULL,
// 			children INTEGER NOT NULL,
// 			babies INTEGER NOT NULL,
// 			adults INTEGER NOT NULL,
// 			milk INTEGER NOT NULL,
// 			diapers INTEGER NOT NULL,
// 			basket INTEGER NOT NULL,
// 			clothing INTEGER NOT NULL,
// 			drugs INTEGER NOT NULL,
// 			other TEXT,
// 			taken BOOLEAN NOT NULL DEFAULT 0
// 		)`,
// 		`CREATE TABLE IF NOT EXISTS logs (
// 			id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			familyID INTEGER NOT NULL,
// 			userID INTEGER NOT NULL,
// 			changeDescription TEXT NOT NULL,
// 			timestamp TEXT NOT NULL,
// 			FOREIGN KEY (familyID) REFERENCES families (id),
// 			FOREIGN KEY (userID) REFERENCES users (id)
// 		)`,
// 		`CREATE TABLE IF NOT EXISTS active_sessions (
// 			user_id INTEGER PRIMARY KEY,
// 			access_token TEXT,
// 			login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 			token_creation_time TIMESTAMP,
// 			FOREIGN KEY(user_id) REFERENCES users(id)
// 		)`,
// 		`CREATE TABLE IF NOT EXISTS revoked_tokens (
// 			id INTEGER PRIMARY KEY AUTOINCREMENT,
// 			jti TEXT NOT NULL,
// 			revoked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 		)`,
// 	}

// 	for _, query := range createTables {
// 		_, err := DB.Exec(query)
// 		if err != nil {
// 			log.Fatal("Failed to execute query:", err)
// 		}
// 	}

// 	// Insert default admin user if not exists
// 	_, err = DB.Exec(`INSERT INTO users (username, password, isAdmin)
// 		SELECT 'admin', 'admin123', 1 WHERE NOT EXISTS
// 		(SELECT 1 FROM users WHERE username = 'admin')`)
// 	if err != nil {
// 		log.Fatal("Failed to insert admin user:", err)
// 	}

// 	fmt.Println("Database initialized successfully!")

// }

package db

import (
	// "context"
	// "database/sql"
	// "fmt"
	"log"
	"os"
	"strconv"

	// "time"

	//
	// _ "github.com/mattn/go-sqlite3"
	// _ "github.com/lib/pq" // PostgreSQL driver
	sqlitecloud "github.com/sqlitecloud/sqlitecloud-go"
)

// var DB *sql.DB

// func InitDB() {
// 	var err error
// 	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
// 		"cnst4x7hhz.g2.sqlite.cloud", // host
// 		8860,                         // port
// 		"aid_app.db",                 // user (database name acts as username in some SQLite Cloud configs)
// 		"etxRvv4Mmrh6nXNddchOveOm1vAP7cwp2UMZWMxgVGw", // password (using API key as password)
// 		"aid_app.db") // dbname

// 	DB, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Failed to connect to SQLite Cloud:", err)
// 	}

// 	// Configure connection pool
// 	DB.SetMaxOpenConns(25)
// 	DB.SetMaxIdleConns(25)
// 	DB.SetConnMaxLifetime(5 * time.Minute)

// 	// Verify connection
// 	if err = DB.Ping(); err != nil {
// 		log.Fatal("Failed to ping database:", err)
// 	}

// 	log.Println("Successfully connected to SQLite Cloud via PostgreSQL protocol!")
// 	createTables()

// 	// dbPath := os.Getenv("DATABASE_URL")
// 	// dbPath := "C:\\Users\\Mad\\Desktop\\help_app_golang\\golang_api_aidapp\\aid_app.db"
// 	// if dbPath == "" {
// 	// dbPath := "/app/aid_app.db"
// 	// dbPath = "C:\\Users\\Mad\\Desktop\\help_app_golang\\golang_api_aidapp\\aid_app.db"
// 	// }

// 	// DB, err = sql.Open("sqlite3", dbPath)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// // Verify connection
// 	// if err = DB.Ping(); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// createTables()
// }

// func InitDB() {
// 	var err error

// 	// Use environment variables for security (set these in Railway)
// 	host := os.Getenv("DB_HOST")
// 	port := os.Getenv("DB_PORT")
// 	user := os.Getenv("DB_USER")
// 	password := os.Getenv("DB_PASSWORD")
// 	dbname := os.Getenv("DB_NAME")

// 	// Convert port to integer
// 	portInt, err := strconv.Atoi(port)
// 	if err != nil {
// 		log.Fatal("Invalid DB_PORT:", err)
// 	}

// 	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require connect_timeout=5",
// 		host,
// 		portInt,
// 		user,
// 		password,
// 		dbname)

// 	// Add connection timeout context
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	DB, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Failed to initialize database connection:", err)
// 	}

// 	// Configure connection pool
// 	DB.SetMaxOpenConns(10) // Reduced for SQLite Cloud's likely limits
// 	DB.SetMaxIdleConns(5)
// 	DB.SetConnMaxLifetime(30 * time.Minute)

// 	// Verify connection with timeout
// 	if err = DB.PingContext(ctx); err != nil {
// 		log.Fatal("Failed to ping database:", err)
// 	}

// 	log.Println("Successfully connected to SQLite Cloud via PostgreSQL protocol!")
// 	createTables()
// }

// Helper struct for DB configuration
type dbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

var DB *sqlitecloud.SQCloud

func InitDB() {
	// Configuration - use environment variables in production
	// config := sqlitecloud.SQCloudConfig{
	// 	Host:    getEnv("SQLITECLOUD_HOST", "cnst4x7hhz.g2.sqlite.cloud"),
	// 	Port:    getEnvInt("SQLITECLOUD_PORT", 8860),
	// 	ApiKey:  getEnv("SQLITECLOUD_APIKEY", "etxRvv4Mmrh6nXNddchOveOm1vAP7cwp2UMZWMxgVGw"),
	// 	Database: getEnv("SQLITECLOUD_DB", "aid_app.db"),
	// 	Timeout: 10 * time.Second,
	// }

	var err error
	DB, err = sqlitecloud.Connect("sqlitecloud://cnst4x7hhz.g2.sqlite.cloud:8860/aid_app.db?apikey=etxRvv4Mmrh6nXNddchOveOm1vAP7cwp2UMZWMxgVGw")
	if err != nil {
		log.Fatalf("❌ Connection failed: %v", err)
	}

	// Verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ Connection verification failed: %v", err)
	}

	log.Println("✅ Successfully connected to SQLite Cloud!")
	createTables()
}

// Helper functions
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func createTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			isAdmin BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS families (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			fullName TEXT NOT NULL,
			nationalID TEXT NOT NULL,
			familyBookID TEXT NOT NULL,
			phoneNumber TEXT NOT NULL,
			familyMembers INTEGER NOT NULL,
			children INTEGER NOT NULL,
			babies INTEGER NOT NULL,
			adults INTEGER NOT NULL,
			milk INTEGER DEFAULT 0,
			diapers INTEGER DEFAULT 0,
			basket INTEGER DEFAULT 0,
			clothing INTEGER DEFAULT 0,
			drugs INTEGER DEFAULT 0,
			other TEXT DEFAULT '',
			taken BOOLEAN DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS active_sessions (
			user_id INTEGER NOT NULL,
			access_token TEXT NOT NULL,
			login_time DATETIME NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS revoked_tokens (
			jti TEXT PRIMARY KEY
		)`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			familyID INTEGER NOT NULL,
			userID INTEGER NOT NULL,
			changeDescription TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			FOREIGN KEY(familyID) REFERENCES families(id),
			FOREIGN KEY(userID) REFERENCES users(id)
		)`,
	}

	for _, table := range tables {
		if err := DB.Execute(table); err != nil {
			log.Fatalf("Failed to execute query: %v\nQuery: %s", err, table)
		}
	}
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Println("Error closing database:", err)
	}
}
