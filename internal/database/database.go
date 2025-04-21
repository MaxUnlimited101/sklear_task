// internal/database/database.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"maxcool.com/weatherapp/internal/models"
)

// DB struct holds the database connection pool
type DB struct {
	SQL *sql.DB
}

// NewDB initializes and returns a new DB instance with the connection pool
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("unable to open database connection: %w", err)
	}

	// Ping the database to verify the connection is established
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DB{SQL: db}, nil
}

// Close closes the database connection pool
func (d *DB) Close() {
	if d.SQL != nil {
		d.SQL.Close()
	}
}

// EnsureDatabaseExists connects to a default database (like 'postgres')
// and executes CREATE DATABASE IF NOT EXISTS for the target database.
// connectionString must be a URL that can connect to the server with create privileges,
// the database name part of this string will be ignored and replaced,
// typically connect to the 'postgres' or an empty database.
func EnsureDatabaseExists(dbName string, connectionString string) error {
	// Parse the connection string URL
	u, err := url.Parse(connectionString)
	if err != nil {
		return fmt.Errorf("invalid connection string URL: %w", err)
	}

	u.Path = "/postgres" // Connect to the default postgres database

	// Reconstruct the connection string for the default database
	connectURL := u.String()
	log.Printf("Attempting to connect to default DB: %s", connectURL)

	// Connect to the default database
	db, err := sql.Open("pgx", connectURL)
	if err != nil {
		return fmt.Errorf("unable to connect to default database (%s): %w", connectURL, err)
	}
	defer db.Close()

	// Ping to verify the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("ping to default database failed: %w", err)
	}
	log.Println("Successfully connected to default database.")

	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)

	log.Printf("Executing: %s", createDBSQL)
	_, err = db.Exec(createDBSQL)
	if err != nil {
		return fmt.Errorf("failed to execute CREATE DATABASE: %w", err)
	}

	log.Printf("Database '%s' ensured to exist (created or already present).", dbName)

	return nil // Success
}

// --- Database Operation Methods (repository layer) ---

// CreateUser inserts a new user into the database
// Returns the ID of the newly created user
func (d *DB) CreateUser(user *models.User) (int, error) {
	var userID int
	err := d.SQL.QueryRow(
		"INSERT INTO users (password, email) VALUES ($1, $2) RETURNING id",
		user.Password, user.Email,
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	user.Id = userID
	return userID, nil
}

// GetUserByEmail retrieves a user by their email address
// Returns the user if found, or nil if not found
func (d *DB) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{} // Create an empty user struct
	err := d.SQL.QueryRow(
		"SELECT id, password, email FROM users WHERE email = $1",
		email,
	).Scan(&user.Id, &user.Password, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// CreateSubscription inserts a new subscription into the database
// Returns the ID of the newly created subscription
func (d *DB) CreateSubscription(sub *models.Subscription) (int, error) {
	var subID int
	err := d.SQL.QueryRow(
		"INSERT INTO subscriptions (user_id, city, condition) VALUES ($1, $2, $3) RETURNING id",
		sub.UserId, sub.City, sub.Condition,
	).Scan(&subID)

	if err != nil {
		return 0, fmt.Errorf("failed to create subscription: %w", err)
	}
	sub.Id = subID
	return subID, nil
}

// GetSubscriptionsByUserID retrieves all subscriptions for a given user ID
// Returns a slice of subscriptions or an error if the query fails
func (d *DB) GetSubscriptionsByUserID(userID int) ([]models.Subscription, error) {
	rows, err := d.SQL.Query("SELECT id, user_id, city, condition FROM subscriptions WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions for user %d: %w", userID, err)
	}
	defer rows.Close()

	subscriptions := []models.Subscription{}

	for rows.Next() {
		var sub models.Subscription

		if err := rows.Scan(&sub.Id, &sub.UserId, &sub.City, &sub.Condition); err != nil {
			log.Printf("Error scanning subscription row for user %d: %v", userID, err) // Log the error but try to continue
			continue                                                                   // Skip this row
		}
		subscriptions = append(subscriptions, sub)
	}

	// Check for errors encountered during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during subscriptions iteration for user %d: %w", userID, err)
	}

	return subscriptions, nil
}

// GetUserByID retrieves a user by their ID
// Returns the user if found, or nil if not found
func (d *DB) GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	err := d.SQL.QueryRow(
		"SELECT id, password, email FROM users WHERE id = $1",
		userID,
	).Scan(&user.Id, &user.Password, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

// UpdateUser updates an existing user in the database
// Returns an error if the update fails
func (d *DB) UpdateUser(user *models.User) error {
	_, err := d.SQL.Exec(
		"UPDATE users SET password = $1, email = $2 WHERE id = $3",
		user.Password, user.Email, user.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user from the database
// Returns an error if the deletion fails
func (d *DB) DeleteUser(userID int) error {
	_, err := d.SQL.Exec(
		"DELETE FROM users WHERE id = $1",
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetSubscriptionByID retrieves a subscription by its ID
// Returns the subscription if found, or nil if not found
// Returns an error if the query fails
func (d *DB) GetSubscriptionByID(subID int) (*models.Subscription, error) {
	sub := &models.Subscription{}
	err := d.SQL.QueryRow(
		"SELECT id, user_id, city, condition FROM subscriptions WHERE id = $1",
		subID,
	).Scan(&sub.Id, &sub.UserId, &sub.City, &sub.Condition)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription by ID: %w", err)
	}
	return sub, nil
}

// UpdateSubscription updates an existing subscription in the database
// Returns an error if the update fails
func (d *DB) UpdateSubscription(sub *models.Subscription) error {
	_, err := d.SQL.Exec(
		"UPDATE subscriptions SET user_id = $1, city = $2, condition = $3 WHERE id = $4",
		sub.UserId, sub.City, sub.Condition, sub.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

// DeleteSubscription deletes a subscription from the database
// Returns an error if the deletion fails
func (d *DB) DeleteSubscription(subID int) error {
	_, err := d.SQL.Exec(
		"DELETE FROM subscriptions WHERE id = $1",
		subID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}
