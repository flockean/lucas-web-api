package database

import (
	"LucasApi/api/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// DB wraps the sql.DB connection
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(cfg *config.DatabaseConfig) (*DB, error) {
	db, err := sql.Open("postgres", cfg.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Initialize database tables
	if err := initializeDatabase(db); err != nil {
		log.Printf("Warning: Error initializing database: %v", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

// Ping tests the database connection
func (db *DB) Ping() error {
	return db.DB.Ping()
}

// GetDBStats returns database statistics
func (db *DB) GetDBStats() sql.DBStats {
	return db.Stats()
}

// initializeDatabase runs the initial SQL schema
func initializeDatabase(db *sql.DB) error {
	sqlStmt, err := os.ReadFile("./db/init.sql")
	if err != nil {
		return fmt.Errorf("could not read init.sql: %v", err)
	}

	_, err = db.Exec(string(sqlStmt))
	if err != nil {
		return fmt.Errorf("error executing init.sql: %v", err)
	}

	log.Println("Database tables initialized successfully")
	return nil
}

// Transaction executes a function within a database transaction
func (db *DB) Transaction(fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Failed to rollback transaction during panic: %v", rbErr)
			}
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Failed to rollback transaction: %v", rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// Health checks database health
func (db *DB) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
