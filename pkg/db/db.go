package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
	"github.com/userblog/management/pkg/config"
)

// connect initializes the database connection (private function)
func Connect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file, using default values")
	}

	dbType := config.GetOrDefaultString("DB_TYPE", "sqlite")
	var connectionString string
	var dialect string

	switch dbType {
	case "postgres":
		host := config.GetOrDefaultString("DB_HOST", "localhost")
		port := config.GetOrDefaultString("DB_PORT", "5432")
		user := config.GetOrDefaultString("DB_USER", "postgres")
		password := config.GetOrDefaultString("DB_PASSWORD", "postgres")
		dbname := config.GetOrDefaultString("DB_NAME", "userblog")
		dialect = "postgres"
		connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	case "mysql":
		host := config.GetOrDefaultString("DB_HOST", "localhost")
		port := config.GetOrDefaultString("DB_PORT", "3306")
		user := config.GetOrDefaultString("DB_USER", "santosh")
		password := config.GetOrDefaultString("DB_PASSWORD", "password")
		dbname := config.GetOrDefaultString("DB_NAME", "userblog")
		dialect = "mysql"
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbname)
	default:
		dialect = "sqlite3"
		connectionString = "./userblog.db"
	}

	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	return db
}
