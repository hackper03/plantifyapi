package dbase

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Database connection function
func ConnectDB() (*sql.DB, error) {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}

	// Verify the connection is working
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not ping the database: %v", err)
	}

	fmt.Printf("Successfully connected to the database\n")
	return db, nil
}


func CreateTables() {
	fmt.Printf("Creating events table")
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		phone VARCHAR(15),
		address TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS categories (
		category_id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		description TEXT
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS plants (
		plant_id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		category_id INT REFERENCES categories(category_id),
		image_url TEXT,
		rating float,
		price float,
		original_price float,
		badgeText VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS plant_inventory (
		inventory_id SERIAL PRIMARY KEY,
		plant_id INT REFERENCES plants(plant_id),
		price DECIMAL(10, 2) NOT NULL,
		stock_quantity INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS services (
		service_id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS cart (
		cart_id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);	
		`,
		`
		CREATE TABLE IF NOT EXISTS cart_items (
		cart_item_id SERIAL PRIMARY KEY,
		cart_id INT REFERENCES cart(cart_id) ON DELETE CASCADE,
		plant_id INT REFERENCES plants(plant_id),
		service_id INT REFERENCES services(service_id),
		quantity INT NOT NULL,
		price DECIMAL(10, 2),
		total DECIMAL(10, 2) GENERATED ALWAYS AS (quantity * price) STORED
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS orders (
		order_id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
		total_price DECIMAL(10, 2) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'Pending',
		shipping_address TEXT,
		order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS order_items (
		order_item_id SERIAL PRIMARY KEY,
		order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
		plant_id INT REFERENCES plants(plant_id),
		service_id INT REFERENCES services(service_id),
		quantity INT,
		price DECIMAL(10, 2),
		total DECIMAL(10, 2) GENERATED ALWAYS AS (quantity * price) STORED
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS payments (
		payment_id SERIAL PRIMARY KEY,
		order_id INT REFERENCES orders(order_id),
		payment_method VARCHAR(50),
		payment_status VARCHAR(50) DEFAULT 'Pending',
		payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		amount DECIMAL(10, 2) NOT NULL
		);
		`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			fmt.Printf("Error executing query: %v\n", err)
			panic("Could not execute query")
		}
		fmt.Printf("Successfully executed query\n")
	}
}