package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/app/models"
	"os"
)

var (
	Db *gorm.DB
)

func ConnectDB() (*gorm.DB, error) {

	/*fmt.Print(os.Getenv("POSTGRES_PASSWORD_FILE"))
	dat, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	if err != nil {
		panic("failed to read password file")*/

	port := "port=" + os.Getenv("DATABASE_PORT")
	host := "host=" + os.Getenv("DATABASE_HOST")
	userq := "user=" + os.Getenv("DATABASE_USER")
	//password := "password=" + strings.TrimSuffix(string(dat), "\n")
	password := "password=postgres"
	name := "database=" + os.Getenv("DATABASE_DB")

	dsn := port + " " + host + " " + userq + " " + password + " " + name + " sslmode=disable TimeZone=America/New_York"
	fmt.Print(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	// Migrate the schema

	if db.AutoMigrate(&models.User{}) != nil {
		fmt.Print("failed to migrate users")
		return db, err
	}
	if db.AutoMigrate(&models.Task{}) != nil {
		fmt.Print("failed to migrate tasks")
		return db, err
	}

	return db, err
}
