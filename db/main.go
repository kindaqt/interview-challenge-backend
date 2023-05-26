package db

import (
	"database/sql"
	"fmt"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Database struct {
	db     *sql.DB
	config Config
}

func NewDatabase(config *Config) (*Database, error) {
	return &Database{
		config: func() Config {
			if config != nil {
				return *config
			}

			return Config{
				Host:     os.Getenv("POSTGRES_HOST"),
				Port:     os.Getenv("POSTGRES_PORT"),
				User:     os.Getenv("POSTGRES_USER"),
				Password: os.Getenv("POSTGRES_PASSWORD"),
				DBName:   os.Getenv("POSTGRES_DB"),
			}
		}(),
	}, nil
}

func (d *Database) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.config.Host, d.config.Port, d.config.User, d.config.Password, d.config.DBName,
	))
	if err != nil {
		return nil, err
	}

	return db, nil
}
