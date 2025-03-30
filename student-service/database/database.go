package database

import (
	"encoding/json"
	"os"

	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func InitDatabade(db *gorm.DB) *Database {
	return &Database{DB: db}
}

type DatabaseConfig struct {
	ConnString    string `json:"connection_string"`
	SingularTable bool   `json:"singular_table"`
}

func ReadDBConfigFromFile(filename string) (*DatabaseConfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var config DatabaseConfig
	decoder := json.NewDecoder(f)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
