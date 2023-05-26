package db

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Runs before every test
func TestMain(m *testing.M) {
	// Get Config
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestQueryRentalByID(t *testing.T) {
	// Start new DB
	d, err := NewDatabase(nil)
	assert.NoError(t, err)

	// Connect to DB
	conn, err := d.Connect()
	assert.NoError(t, err)
	defer conn.Close()

	// Start a transaction
	tx, err := conn.Begin()
	assert.NoError(t, err)

	// Query DB for rental by id
	rental, err := d.QueryRentalByID(tx, 1)
	assert.NoError(t, err)

	t.Logf("%d %T", rental.UserID, rental.UserID)

	// Query DB for user by id
	user, err := d.QueryUserByID(tx, rental.UserID)
	assert.NoError(t, err)

	t.Log(rental, user)
}

func TestQueryRentals(t *testing.T) {
	// Start new DB
	d, err := NewDatabase(nil)
	assert.NoError(t, err)

	// Connect to DB
	conn, err := d.Connect()
	assert.NoError(t, err)
	defer conn.Close()

	// Start a transaction
	tx, err := conn.Begin()
	assert.NoError(t, err)

	// Query rentals
	priceMin := 9000
	priceMax := 75000
	queryParams := QueryRentalsParams{
		PriceMin: &priceMin,
		PriceMax: &priceMax,
	}

	rentals, err := d.QueryRentals(tx, &queryParams)
	assert.NoError(t, err)

	jsonData, _ := json.MarshalIndent(rentals, "", " ")
	t.Log(string(jsonData))
}
