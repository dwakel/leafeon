package main

import (
	"fmt"
	"leafeon/migrator"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	migrationType := os.Args[0]
	argsConnString := os.Args[2]
	argsPath := os.Args[1]

	dsn := fmt.Sprintf(argsConnString)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}
	migrate := migrator.New(database, argsPath)
	if migrationType == "up" {
		migErr := migrate.Up()
		if migErr != nil {
			fmt.Println(migErr)
		}
	} else if migrationType == "down" {
		migErr := migrate.Up()
		if migErr != nil {
			fmt.Println(migErr)
		}
	} else {
		fmt.Println("Invalid command line args: ", migrationType)
	}
}
