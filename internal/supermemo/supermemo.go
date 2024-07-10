package supermemo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/takacs/donkey/internal/database"
	"github.com/takacs/donkey/internal/review"
)

type Supermemo struct {
	ID             uint
	CardID         uint
	Repetition     int
	EasinessFactor float64
	Interval       int
	NextReview     time.Time
}

type SupermemoDb struct {
	db *sql.DB
}

func (c *SupermemoDb) Insert(cardId uint) error {
	_, err := c.db.Exec(
		"INSERT INTO supermemo(card_id, repetition, easiness_factor, interval, next_review_time) VALUES( ?, ?, ?, ?, ?)",
		cardId,
		0,
		2.5,
		0,
		time.Time{},
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *SupermemoDb) Delete(id uint) error {
	_, err := c.db.Exec("DELETE FROM supermemo WHERE id = ?", id)
	return err
}

func New() (*SupermemoDb, error) {
	db, err := database.OpenDb()
	if err != nil {
		return nil, errors.New("couldn't open db")
	}
	supermemoDb := SupermemoDb{db: db}
	return &supermemoDb, nil
}

func (c *SupermemoDb) Close() error {
	if err := c.db.Close(); err != nil {
		return err
	}
	return nil
}

func (c *SupermemoDb) GetCardsSupermemo(cardId uint) (Supermemo, error) {
	cardIdStr := strconv.Itoa(int(cardId))
	query := fmt.Sprintf("SELECT * FROM supermemo WHERE card_id = " + cardIdStr)
	rows, err := c.db.Query(query)
	if err != nil {
		return Supermemo{}, err
	}

	if !rows.Next() {
		log.Printf("inserting card id %v", cardId)
		err := c.Insert(cardId)
		if err != nil {
			return Supermemo{}, err
		}
		rows, err = c.db.Query(query)
		if err != nil {
			return Supermemo{}, err
		}
	}

	err = rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	// TODO just look at it
	query = fmt.Sprintf("SELECT * FROM supermemo WHERE card_id = " + cardIdStr)
	rows, err = c.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var supermemo Supermemo
	for rows.Next() {
		err = rows.Scan(
			&supermemo.ID,
			&supermemo.CardID,
			&supermemo.Repetition,
			&supermemo.EasinessFactor,
			&supermemo.Interval,
			&supermemo.NextReview,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("returned %v", supermemo)
	return supermemo, nil
}

func (c *SupermemoDb) GetXSoonestReviewTimeCardIds(x int) ([]uint, error) {
	xStr := strconv.Itoa(x)
	query := fmt.Sprintf("SELECT card_id FROM supermemo ORDER BY next_review_time LIMIT " + xStr)
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}

	var cardId uint
	cardIds := []uint{}
	for rows.Next() {
		err = rows.Scan(
			&cardId,
		)
		if err != nil {
			return nil, err
		}
		cardIds = append(cardIds, cardId)
	}

	return cardIds, nil
}

func (s SupermemoDb) updateSupermemo(supermemoId uint, n, I int, EF float64) error {
	interval := time.Hour * 24 * time.Duration(I)
	nextReviewTime := time.Now().Add(interval)
	_, err := s.db.Exec(
		"UPDATE supermemo SET repetition = ?, interval = ?, easiness_factor = ?, next_review_time = ? WHERE id = ?",
		n,
		I,
		EF,
		nextReviewTime,
		supermemoId,
	)
	if err != nil {
		return err
	}
	log.Printf("updated %v", supermemoId)
	return nil
}

func UpdateCardParams(cardId uint, grade review.Grade) error {
	supermemoDb, err := New()
	if err != nil {
		return err
	}
	supermemo, err := supermemoDb.GetCardsSupermemo(cardId)
	if err != nil {
		return err
	}
	log.Printf("calculating new params for card %v\n", cardId)

	var n int
	var I int
	var EF float64
	if grade == review.Again {
		n = 0
		I = 1
	} else {
		if n == 0 {
			I = 1
		} else if n == 1 {
			if grade == review.Hard {
				I = 3
			} else {
				I = 6
			}
		} else {
			I = int(float64(supermemo.Interval) * supermemo.EasinessFactor)
		}
	}
	n++

	gradeValue := int(grade)
	EF = supermemo.EasinessFactor + (0.1 - float64(5-gradeValue)*(0.08+float64(5-gradeValue)*0.02))
	if EF > 1.3 {
		EF = 1.3
	}
	log.Printf("new params: n=%v, I=%v, EF=%v\n", n, I, EF)

	err = supermemoDb.updateSupermemo(supermemo.ID, n, I, EF)
	if err != nil {
		return err
	}
	return nil

}
