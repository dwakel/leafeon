package migrator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"gorm.io/gorm"
)

type Migrator struct {
	database *gorm.DB
	path     string
}

func New(database *gorm.DB, path string) *Migrator {
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
	var cd string
	if m.path == "" {
		cd = "./migrations"
	} else {
		cd = m.path
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
		//fmt.Println(ext)
		data, err := ioutil.ReadFile(cd + "/" + file.Name())
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
	return nil
}

func getDoubleExt(path string) (string, error) {
	arr := strings.Split(path, ".")
	if len(arr) <= 2 {
		return "", errors.New("Invalid File")
	}
	return strings.Join(arr[len(arr)-2:], "."), nil
}
