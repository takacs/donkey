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

// Opens donkey db with all neccessary tables created
// If donkey db doesn't exist it this function creates it as well
func OpenDb() (*sql.DB, error) {
	path, err := GetDbPath()
	if err != nil {
		return nil, err
	}

	db, err := InitDatabase(path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Return datadir containing donkey db
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

func InitDatabase(path string) (*sql.DB, error) {
	dbPath := filepath.Join(path, donkeyDbName+".db")
	_, err := os.Stat(dbPath)

	if err == nil {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal("can't open db")
		}
		return db, nil
	}

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
