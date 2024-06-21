package card

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/takacs/donkey/internal/database"
)

func TestGetCard(t *testing.T) {
	tests := []struct {
		want Card
	}{
		{
			want: Card{
				ID:    1,
				Front: "get milk",
				Back:  "groceries",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.want.Front, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if _, err := tDB.Insert(tc.want.Front, tc.want.Back, tc.want.Deck); err != nil {
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

func setup() *CardDb {
	path := filepath.Join(os.TempDir())
	db, err := database.InitDatabase(path)
	if err != nil {
		log.Fatal(err)
	}
	t := CardDb{db}
	return &t
}

func teardown(tDB *CardDb) {
	tDB.Close()
	path := filepath.Join(os.TempDir(), "donkey.db")
	os.Remove(path)

}
