package main

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"time"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"todo", "in progress", "done"}[s]
}

type card struct {
	ID      uint
	Front   string
	Back    string
	Deck    string
	Status  string
	Created time.Time
}

func (c card) FilterValue() string {
	return c.Front
}

func (c card) Title() string {
	return c.Front
}

func (c card) Description() string {
	return c.Back
}

func (s status) Next() int {
	if s == done {
		return int(todo)
	}
	return int(s + 1)
}

func (s status) Prev() int {
	if s == todo {
		return int(done)
	}
	return int(s - 1)
}

func (s status) Int() int {
	return int(s)
}

type cardDB struct {
	db      *sql.DB
	dataDir string
}

func initCardDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

func (c *cardDB) tableExists(name string) bool {
	if _, err := c.db.Query("SELECT * FROM cards"); err == nil {
		return true
	}
	return false
}

func (c *cardDB) createTable() error {
	_, err := c.db.Exec(`CREATE TABLE "cards" ( "id" INTEGER, "front" TEXT NOT NULL, "back" TEXT, "deck" TEXT, "status" TEXT, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}

func (c *cardDB) insert(front, back, deck string) error {
	_, err := c.db.Exec(
		"INSERT INTO cards(front, back, deck, status, created) VALUES( ?, ?, ?, ?, ?)",
		front,
		back,
		deck,
		todo.String(),
		time.Now())
	return err
}

func (c *cardDB) delete(id uint) error {
	_, err := c.db.Exec("DELETE FROM cards WHERE id = ?", id)
	return err
}

func (c *cardDB) update(card card) error {
	// Get the existing state of the card we want to update.
	orig, err := c.getcard(card.ID)
	if err != nil {
		return err
	}
	orig.merge(card)
	_, err = c.db.Exec(
		"UPDATE cards SET name = ?, project = ?, deck = ?, status = ? WHERE id = ?",
		orig.Front,
		orig.Back,
		orig.Deck,
		orig.Status,
		orig.ID)
	return err
}

func (orig *card) merge(t card) {
	uValues := reflect.ValueOf(&t).Elem()
	oValues := reflect.ValueOf(orig).Elem()
	for i := 0; i < uValues.NumField(); i++ {
		uField := uValues.Field(i).Interface()
		if oValues.CanSet() {
			if v, ok := uField.(int64); ok && uField != 0 {
				oValues.Field(i).SetInt(v)
			}
			if v, ok := uField.(string); ok && uField != "" {
				oValues.Field(i).SetString(v)
			}
		}
	}
}

func (c *cardDB) getcards() ([]card, error) {
	var cards []card
	rows, err := c.db.Query("SELECT * FROM cards")
	if err != nil {
		return cards, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var card card
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

func (c *cardDB) getcardsByStatus(status string) ([]card, error) {
	var cards []card
	rows, err := c.db.Query("SELECT * FROM cards WHERE status = ?", status)
	if err != nil {
		return cards, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var card card
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

func (c *cardDB) getcard(id uint) (card, error) {
	var card card
	err := c.db.QueryRow("SELECT * FROM cards WHERE id = ?", id).
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
