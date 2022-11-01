package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Migrator struct {
	database *gorm.DB
	path     string
}

type Migration struct {
	Migration string    `pg:"migration"`
	CreatedAt time.Time `pg:"created_at"`
}

const track = "_migrations"

const trackSchema = `CREATE TABLE IF NOT EXISTS public._migrations (
    migration text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);`

type Migrators interface {
	Up() error
	Down() error
}

func New(database *gorm.DB, path string) Migrators {
	return &Migrator{database, path}
}

func (m *Migrator) Up() error {
	err := migrate(m, "up")
	if err != nil {
		return err
	}
	return nil
}

func (m *Migrator) Down() error {
	err := migrate(m, "down")
	if err != nil {
		return err
	}
	return nil
}

func migrate(m *Migrator, mType string) error {
	var cd string = m.path
	if cd == "" {
		cd = "migrations"
	}
	files, err := ioutil.ReadDir(cd)
	if err != nil {
		return errors.New(fmt.Sprintf("Migration failed: %v", err.Error()))
	}
	for _, file := range files {
		ext, err := getDoubleExt(file.Name())
		if err != nil {
			return errors.New(fmt.Sprintf("Migration failed: %v", err.Error()))
		}
		if file.IsDir() || ext != mType+".sql" {
			continue
		}
		fmt.Println(file.Name())

		//Check if Migration was already run
		m.database.Exec(trackSchema)
		var record Migration
		upFName, _ := getWithoutDoubleExt(file.Name())
		raw := "SELECT * FROM public._migrations WHERE migration = @name"
		_ = m.database.Raw(raw, sql.Named("name", upFName+".up.sql")).Scan(&record).Error
		if mType == "up" {
			if record.Migration != file.Name() {
				data, err := ioutil.ReadFile(filepath.Join(cd, file.Name()))
				if err != nil {
					fmt.Println()
					return errors.New(fmt.Sprintf("Migration failed: %v", err.Error()))
				}
				err = m.database.Exec(string(data)).Error
				if err != nil {
					fmt.Println()
					return errors.New(fmt.Sprintf("Migration failed: %v.\n Error: %v", file.Name(), err.Error()))
				}
				_ = m.database.Exec(fmt.Sprintf("INSERT INTO public._migrations (migration) VALUES ('%v')", file.Name()))
				fmt.Println(" Done!")
			} else if err == nil {
				fmt.Println(" Skipped!")
			} else {
				return errors.New(fmt.Sprintf("Migration failed: %v", err.Error()))
			}
		}
		if mType == "down" {
			if record.Migration != upFName+".up.sql" {
				fmt.Println(" Skipped!")
			} else if err == nil && record.Migration == upFName+".up.sql" {
				data, err := ioutil.ReadFile(filepath.Join(cd, file.Name()))
				if err != nil {
					fmt.Println()
					return errors.New(fmt.Sprintf("Migration failed: %v", err.Error()))
				}
				err = m.database.Exec(string(data)).Error
				if err != nil {
					fmt.Println()
					return errors.New(fmt.Sprintf("Migration failed: %v.\n Error: %v", file.Name(), err.Error()))
				}
				upName, _ := getWithoutDoubleExt(file.Name())
				_ = m.database.Exec(fmt.Sprintf("DELETE FROM public._migrations WHERE migration = '%v'", upName+".up.sql"))
				fmt.Println(" Done!")

			} else {
				return errors.New(fmt.Sprintf("Migration failed: %v", err.Error()))
			}
		}
	}
	fmt.Println(" Executed all migrations!")
	return nil
}

func getDoubleExt(path string) (string, error) {
	arr := strings.Split(path, ".")
	if len(arr) <= 2 {
		return "", errors.New("Invalid File")
	}
	return strings.Join(arr[len(arr)-2:], "."), nil
}

func getWithoutDoubleExt(path string) (string, error) {
	arr := strings.Split(path, ".")
	if len(arr) <= 2 {
		return "", errors.New("Invalid File")
	}
	return strings.Join(arr[:len(arr)-2], "."), nil
}
