package card

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/takacs/donkey/internal/database"
	"log"
	"time"
)

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

func (s Status) String() string {
	return [...]string{"todo", "in progress", "done"}[s]
}

type Card struct {
	ID      uint
	Front   string
	Back    string
	Deck    string
	Status  string
	Created time.Time
}

type CardDb struct {
	db *sql.DB
}

func (c *CardDb) Insert(front, back, deck string) error {
	_, err := c.db.Exec(
		"INSERT INTO card(front, back, deck, status, created) VALUES( ?, ?, ?, ?, ?)",
		front,
		back,
		deck,
		Todo.String(),
		time.Now())
	return err
}

func (c *CardDb) Delete(id uint) error {
	_, err := c.db.Exec("DELETE FROM card WHERE id = ?", id)
	return err
}

func New() (*CardDb, error) {
	db, err := database.OpenDb()
	if err != nil {
		log.Fatal("couldn't open db")
		return nil, errors.New("couldn't open db")
	}
	cardDb := CardDb{db: db}
	if !cardDb.tableExists() {
		err := cardDb.createTable()
		if err != nil {
			return nil, err
		}
	}
	return &cardDb, nil
}

func (c *CardDb) Close() error {
	if err := c.db.Close(); err != nil {
		log.Fatal("failed closing db")
		return errors.New("failed closing db")
	}
	return nil
}

func (c *CardDb) GetCards() ([]Card, error) {
	var cards []Card
	rows, err := c.db.Query("SELECT * FROM card")
	if err != nil {
		return cards, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var card Card
		err = rows.Scan(
			&card.ID,
			&card.Front,
			&card.Back,
			&card.Deck,
			&card.Status,
			&card.Created,
		)
		if err != nil {
			return cards, err
		}
		cards = append(cards, card)
	}
	return cards, err
}

func (c *CardDb) GetCardsByStatus(status string) ([]Card, error) {
	var cards []Card
	rows, err := c.db.Query("SELECT * FROM card WHERE status = ?", status)
	if err != nil {
		return cards, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var card Card
		err = rows.Scan(
			&card.ID,
			&card.Front,
			&card.Back,
			&card.Deck,
			&card.Status,
			&card.Created,
		)
		if err != nil {
			return cards, err
		}
		cards = append(cards, card)
	}
	return cards, err
}

func (c *CardDb) GetCard(id uint) (Card, error) {
	var card Card
	err := c.db.QueryRow("SELECT * FROM card WHERE id = ?", id).
		Scan(
			&card.ID,
			&card.Front,
			&card.Back,
			&card.Deck,
			&card.Status,
			&card.Created,
		)
	return card, err
}

func (c *CardDb) tableExists() bool {
	if _, err := c.db.Query("SELECT * FROM card"); err == nil {
		return true
	}
	return false
}

func (c *CardDb) createTable() error {
	_, err := c.db.Exec(`CREATE TABLE "card" ( "id" INTEGER, "front" TEXT NOT NULL, "back" TEXT, "deck" TEXT, "status" TEXT, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}
