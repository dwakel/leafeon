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
	defer catchException()
	migrationType := os.Args[3]
	argsConnString := flag.String("connstr", "", "provide relevant connection string to establish connection to PostgreSQL database")
	argsPath := flag.String("src", "", "provide a path to source directory to peform migrations from")
	flag.Parse()
	if *argsConnString == "" {
		fmt.Println("Please provide relevant connection string: -connstr=")
		return
	}
	if *argsPath == "" {
		fmt.Println("Please provide directory for migrations: -src=")
		return
	}
	database, err := gorm.Open(postgres.Open(fmt.Sprintf(*argsConnString)), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	migrate := migrator.New(database, *argsPath)
	runMigration(migrationType, migrate)
	
}

func runMigration(migrationType string, migrate migrator.Migrators) {
	if migrationType == "up" {
		migErr := migrate.Up()
		if migErr != nil {
			fmt.Println(migErr)
		}
	} else if migrationType == "down" {
		migErr := migrate.Down()
		if migErr != nil {
			fmt.Println(migErr)
		}
	} else {
		fmt.Println("Invalid command line args: ", migrationType)
	}
}

func catchException() {
	if ex := recover(); ex != nil {
		fmt.Println("Unexpected error occured.'")
	}
}
