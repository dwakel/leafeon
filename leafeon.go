package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dwakel/leafeon/migrator"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	migrationType := os.Args[1]

	argsConnString := flag.String("connstr", "", "provide relevant connection string to establish connection to PostgreSQL database")
	argsPath := flag.String("src", "", "provide a path to source directory to peform migrations from")

	if *argsConnString == "" {
		fmt.Println("Please provide relevant connection string: -connstr=")
		return
	}
	if *argsPath == "" {
		fmt.Println("Please provide directory for migrations: -src=")
		return
	}

	dsn := fmt.Sprintf(*argsConnString)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		return
	}
	migrate := migrator.New(database, *argsPath)
	if migrationType == "up" {
		migErr := migrate.Up()
		if migErr != nil {
			fmt.Println(migErr)
			return
		}
	} else if migrationType == "down" {
		migErr := migrate.Up()
		if migErr != nil {
			fmt.Println(migErr)
			return
		}
	} else {
		fmt.Println("Invalid command line args: ", migrationType)
		return 
	}
}
