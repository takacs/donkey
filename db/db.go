package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
	"log"
	"os"
	"path/filepath"
	"reflect"
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

func (c Card) FilterValue() string {
	return c.Front
}

func (c Card) Title() string {
	return c.Front
}

func (c Card) Description() string {
	return c.Back
}

func (s Status) Int() int {
	return int(s)
}

type CardDB struct {
	Db      *sql.DB
	dataDir string
}

func (c *CardDB) TableExists(name string) bool {
	if _, err := c.Db.Query("SELECT * FROM cards"); err == nil {
		return true
	}
	return false
}

func (c *CardDB) CreateTable() error {
	_, err := c.Db.Exec(`CREATE TABLE "cards" ( "id" INTEGER, "front" TEXT NOT NULL, "back" TEXT, "deck" TEXT, "status" TEXT, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}

func (c *CardDB) Insert(front, back, deck string) error {
	_, err := c.Db.Exec(
		"INSERT INTO cards(front, back, deck, status, created) VALUES( ?, ?, ?, ?, ?)",
		front,
		back,
		deck,
		Todo.String(),
		time.Now())
	return err
}

func (c *CardDB) Delete(id uint) error {
	_, err := c.Db.Exec("DELETE FROM cards WHERE id = ?", id)
	return err
}

func (c *CardDB) Update(card Card) error {
	orig, err := c.Getcard(card.ID)
	if err != nil {
		return err
	}
	orig.merge(card)
	_, err = c.Db.Exec(
		"UPDATE cards SET front = ?, back = ?, deck = ?, status = ? WHERE id = ?",
		orig.Front,
		orig.Back,
		orig.Deck,
		orig.Status,
		orig.ID)
	return err
}

func (orig *Card) merge(t Card) {
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

func (c *CardDB) Getcards() ([]Card, error) {
	var cards []Card
	rows, err := c.Db.Query("SELECT * FROM cards")
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

func (c *CardDB) GetcardsByStatus(status string) ([]Card, error) {
	var cards []Card
	rows, err := c.Db.Query("SELECT * FROM cards WHERE status = ?", status)
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

func (c *CardDB) Getcard(id uint) (Card, error) {
	var card Card
	err := c.Db.QueryRow("SELECT * FROM cards WHERE id = ?", id).
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

func OpenDb(path string) (*CardDB, error) {
	db, err := sql.Open("sqlite3", filepath.Join(path, "cards.db"))
	if err != nil {
		return nil, err
	}
	c := CardDB{db, path}
	if !c.TableExists("cards") {
		err := c.CreateTable()
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}

func GetDbPath(app string) (string, error) {
	scope := gap.NewScope(gap.User, app)
	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
		return "", errors.New("cant get datadir")
	}
	// create the app base dir, if it doesn't exist
	var cardDir string
	if len(dirs) > 0 {
		cardDir = dirs[0]
	} else {
		cardDir, _ = os.UserHomeDir()
	}
	if err := initCardDir(cardDir); err != nil {
		log.Fatal(err)
		return "", errors.New("cant init datadir")
	}
	fmt.Println(cardDir)
	return cardDir, nil
}

func initCardDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o750)
		}
		return err
	}
	return nil
}

func SetupPath() string {
	cardDir, err := GetDbPath("cards")
	if err != nil {
		fmt.Println("error getting db path")
	}
	return cardDir
}
