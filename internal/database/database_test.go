package database

import (
	"database/sql"
	gap "github.com/muesli/go-app-paths"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenDb(t *testing.T) {
	_, err := OpenDb()
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDbPath(t *testing.T) {
	path, err := GetDbPath()
	if err != nil {
		t.Fatal(err)
	}
	scope := gap.NewScope(gap.User, donkeyDbName)
	dirs, err := scope.DataDirs()
	if err != nil {
		t.Fatal(err)
	}

	if dirs[0] != path {
		t.Fatal("wrong path")
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestInitDatabase(t *testing.T) {
	path := filepath.Join(os.TempDir())
	db, err := InitDatabase(path)
	if err != nil {
		log.Fatal(err)
	}
	defer teardown(db)

	exists := tableExistsInDatabase(cardTable, db)
	if !exists {
		t.Fatal("card table wasn't created in db")
	}

	exists = tableExistsInDatabase(reviewTable, db)
	if !exists {
		t.Fatal("review table wasn't created in db")
	}

	exists = tableExistsInDatabase(supermemoTable, db)
	if !exists {
		t.Fatal("supermemo table wasn't created in db")
	}
}

func TestInitDataDir(t *testing.T) {
	path := filepath.Join(os.TempDir())
	err := initDataDir(path)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(path)

	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
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
