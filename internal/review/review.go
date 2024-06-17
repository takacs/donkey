package review

import (
	"database/sql"
	"errors"
	"github.com/takacs/donkey/internal/database"
	"log"
	"time"
)

type Grade int

const (
	Easy Grade = iota
	Good
	Hard
	Again
)

type Review struct {
	ID         uint
	CardID     uint
	Grade      Grade
	ReviewTime time.Time
}

type ReviewDb struct {
	db *sql.DB
}

func (c *ReviewDb) Insert(card_id uint, grade Grade) error {
	_, err := c.db.Exec(
		"INSERT INTO review(card_id, grade, reviewed) VALUES( ?, ?, ?)",
		card_id,
		grade,
		time.Now())
	return err
}

func (c *ReviewDb) Delete(id uint) error {
	_, err := c.db.Exec("DELETE FROM review WHERE id = ?", id)
	return err
}

func New() (*ReviewDb, error) {
	db, err := database.OpenDb()
	if err != nil {
		log.Fatal("couldn't open db")
		return nil, errors.New("couldn't open db")
	}
	cardDb := ReviewDb{db: db}
	if !cardDb.tableExists() {
		err := cardDb.createTable()
		if err != nil {
			return nil, err
		}
	}
	return &cardDb, nil
}

func (c *ReviewDb) Close() error {
	if err := c.db.Close(); err != nil {
		log.Fatal("failed closing db")
		return errors.New("failed closing db")
	}
	return nil
}

func (c *ReviewDb) tableExists() bool {
	if _, err := c.db.Query("SELECT * FROM review"); err == nil {
		return true
	}
	return false
}

func (c *ReviewDb) createTable() error {
	_, err := c.db.Exec(`
        CREATE TABLE "review"
        (
        "id" INTEGER,
        "card_id" INTEGER,
        "grade" INTEGER,
        "reviewed" DATETIME,
        PRIMARY KEY("id" AUTOINCREMENT),
        FOREIGN KEY (card_id) REFERENCES card(id))`)
	return err
}
