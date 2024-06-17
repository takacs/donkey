package card

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		want Card
	}{
		{
			want: Card{
				ID:     1,
				Front:  "anki",
				Back:   "Anki is a program which makes remembering things easy.",
				Deck:   "default",
				Status: "todo",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.want.Front, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.Insert(tc.want.Front, tc.want.Back, tc.want.Deck); err != nil {
				t.Fatalf("unable to insert cards: %v", err)
			}
			cards, err := tDB.GetCards()
			if err != nil {
				t.Fatalf("unable to get cards: %v", err)
			}
			tc.want.Created = cards[0].Created
			if !reflect.DeepEqual(tc.want, cards[0]) {
				t.Fatalf("got %v, want %v", tc.want, cards)
			}
			if err := tDB.Delete(1); err != nil {
				t.Fatalf("unable to delete cards: %v", err)
			}
			cards, err = tDB.GetCards()
			if err != nil {
				t.Fatalf("unable to get cards: %v", err)
			}
			if len(cards) != 0 {
				t.Fatalf("expected cards to be empty, got: %v", cards)
			}
		})
	}
}

func TestGetCard(t *testing.T) {
	tests := []struct {
		want Card
	}{
		{
			want: Card{
				ID:     1,
				Front:  "get milk",
				Back:   "groceries",
				Status: Todo.String(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.want.Front, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.Insert(tc.want.Front, tc.want.Back, tc.want.Deck); err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			Card, err := tDB.GetCard(tc.want.ID)
			if err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			tc.want.Created = Card.Created
			if !reflect.DeepEqual(Card, tc.want) {
				t.Fatalf("got: %#v, want: %#v", Card, tc.want)
			}
		})
	}
}

func TestGetCardsByStatus(t *testing.T) {
	tests := []struct {
		want Card
	}{
		{
			want: Card{
				ID:     1,
				Front:  "get milk",
				Back:   "groceries",
				Deck:   "default",
				Status: Todo.String(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.want.Front, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.Insert(tc.want.Front, tc.want.Back, tc.want.Deck); err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			Cards, err := tDB.GetCardsByStatus(tc.want.Status)
			if err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			if len(Cards) < 1 {
				t.Fatalf("expected 1 value, got %#v", Cards)
			}
			tc.want.Created = Cards[0].Created
			if !reflect.DeepEqual(Cards[0], tc.want) {
				t.Fatalf("got: %#v, want: %#v", Cards, tc.want)
			}
		})
	}
}

func setup() *CardDb {
	path := filepath.Join(os.TempDir(), "test.db")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	t := CardDb{db}
	if !t.tableExists() {
		err := t.createTable()
		if err != nil {
			log.Fatal(err)
		}
	}
	return &t
}

func teardown(tDB *CardDb) {
	tDB.Close()
	path := filepath.Join(os.TempDir(), "test.db")
	os.Remove(path)

}
