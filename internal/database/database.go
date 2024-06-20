package database

import (
	"database/sql"
	"fmt"
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
	db, err := InitDatabase(path)
	if err != nil {
		log.Fatal("cant get stat")
	}

	return db, nil

}

func InitDatabase(path string) (*sql.DB, error) {
	dbPath := filepath.Join(path, donkeyDbName+".db")
	_, err := os.Stat(dbPath)

	// no error means db exists
	// we assume it was set up correctly for now
	if err == nil {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("can't open db")
		}
		return db, nil
	}

	// if we get a not exists error we init db
	if os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("could not create")
		}

		if exists := tableExistsInDatabase(cardTable, db); !exists {
			createTable(cardTable, cardSchema, db)
		}

		if exists := tableExistsInDatabase(reviewTable, db); !exists {
			createTable(reviewTable, reviewSchema, db)
		}

		if exists := tableExistsInDatabase(supermemoTable, db); !exists {
			createTable(reviewTable, supermemoSchema, db)
		}

		return db, nil
	}

	// we got an error that is not "not exist" so we return it
	return nil, err
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

func tableExistsInDatabase(tableName string, db *sql.DB) bool {
	if _, err := db.Query(fmt.Sprintf("SELECT * FROM %v", tableName)); err == nil {
		return true
	}
	return false
}

func createTable(name, schema string, db *sql.DB) {
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("couldn't create %v table", name)
	}
}
