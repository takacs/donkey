package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
	"log"
	"os"
	"path/filepath"
)

const (
	donkeyDbName = "donkey"
)

func GetDbPath() (string, error) {
	scope := gap.NewScope(gap.User, donkeyDbName)
	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
	}

	var dir string
	// TODO i did it like this but don't remember why, understand soon
	if len(dirs) > 0 {
		dir = dirs[0]
	} else {
		dir, _ = os.UserHomeDir()
	}

	if err := initDataDir(dir); err != nil {
		log.Fatal(err)
	}
	return dir, nil
}

func OpenDb() (*sql.DB, error) {
	path, err := GetDbPath()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", filepath.Join(path, donkeyDbName+".db"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func initDataDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o750)
		}
		return err
	}
	return nil
}
