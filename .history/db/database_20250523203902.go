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
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	// dbPath := os.Getenv("DATABASE_URL")
	// dbPath := "C:\\Users\\Mad\\Desktop\\help_app_golang\\golang_api_aidapp\\aid_app.db"
	// if dbPath == "" {
	dbPath := "/app/aid_app.db"
	// dbPath = "C:\\Users\\Mad\\Desktop\\help_app_golang\\golang_api_aidapp\\aid_app.db"
	// }

	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Verify connection
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables()
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
		if _, err := DB.Exec(table); err != nil {
			log.Fatal(err)
		}
	}
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Println("Error closing database:", err)
	}
}
