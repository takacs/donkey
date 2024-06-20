package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestInitDatabase(t *testing.T) {
	db := setup()
	defer teardown(db)

	exists := tableExistsInDatabase(cardTable, db)
	if !exists {
		log.Fatal("card table wasn't created in db")
	}

	exists = tableExistsInDatabase(reviewTable, db)
	if !exists {
		log.Fatal("review table wasn't created in db")
	}

	exists = tableExistsInDatabase(supermemoTable, db)
	if !exists {
		log.Fatal("supermemo table wasn't created in db")
	}
}

func TestTableExistsInDatabase(t *testing.T) {
	db := setup()
	defer teardown(db)

	db.Exec("CREATE TABLE test (id INTEGER)")
	if exists := tableExistsInDatabase("test", db); !exists {
		log.Fatal("table test was created but function denies existence")
	}
}

func TestCreateTable(t *testing.T) {
	db := setup()
	defer teardown(db)

	createTable("test", "CREATE TABLE test (id INTEGER)", db)
	if exists := tableExistsInDatabase("test", db); !exists {
		log.Fatal("table test was created but function denies existence")
	}
}

func setup() *sql.DB {
	path := filepath.Join(os.TempDir())
	db, err := InitDatabase(path)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func teardown(tdb *sql.DB) {
	tdb.Close()
	path := filepath.Join(os.TempDir(), "donkey.db")
	os.Remove(path)

}
