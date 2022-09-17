package migrator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type Migrator struct {
	database *gorm.DB
	path     string
}

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
		fmt.Println(" Done!")
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
