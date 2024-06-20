package supermemo

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/takacs/donkey/internal/database"
	"log"
	"time"
)

type Supermemo struct {
	ID             uint
	CardID         uint
	Repetition     uint
	EasinessFactor int
	Interval       time.Time
}

type SupermemoDb struct {
	db *sql.DB
}

func (c *SupermemoDb) Insert(cardId uint) error {
	_, err := c.db.Exec(
		"INSERT INTO card(card_id) VALUES( ?)",
		cardId,
	)
	return err
}

func (c *SupermemoDb) Delete(id uint) error {
	_, err := c.db.Exec("DELETE FROM supermemo WHERE id = ?", id)
	return err
}

func New() (*SupermemoDb, error) {
	db, err := database.OpenDb()
	if err != nil {
		log.Fatal("couldn't open db")
		return nil, errors.New("couldn't open db")
	}
	supermemoDb := SupermemoDb{db: db}
	return &supermemoDb, nil
}

func (c *SupermemoDb) Close() error {
	if err := c.db.Close(); err != nil {
		log.Fatal("failed closing db")
		return errors.New("failed closing db")
	}
	return nil
}

func (c *SupermemoDb) GetCardsSupermemo(cardId uint) Supermemo {
	query := fmt.Sprintf("SELECT * FROM supermemo WHERE card_id = %v", cardId)
	rows, err := c.db.Query(query)
	if err != nil {
		log.Fatalf("no entry for cardid %v in supermemo db", cardId)
	}
	var supermemo Supermemo
	for rows.Next() {
		err = rows.Scan(
			&supermemo.ID,
			&supermemo.CardID,
			&supermemo.Repetition,
			&supermemo.EasinessFactor,
			&supermemo.Interval,
		)
		if err != nil {
			log.Fatal("cant scan query response from supermemo table")
		}
	}
	return supermemo
}
