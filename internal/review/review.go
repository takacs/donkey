package review

import (
	"database/sql"
	"errors"
	"github.com/takacs/donkey/internal/database"
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
		return nil, errors.New("couldn't open db")
	}
	cardDb := ReviewDb{db: db}
	return &cardDb, nil
}

func (c *ReviewDb) Close() error {
	if err := c.db.Close(); err != nil {
		return errors.New("failed closing db")
	}
	return nil
}
