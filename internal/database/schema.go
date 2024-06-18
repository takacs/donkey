package database

var cardTable = "card"

var cardSchema = `
	CREATE TABLE "card"
	(
	"id" INTEGER,
	"front" TEXT NOT NULL,
	"back" TEXT,
	"deck" TEXT,
	"status" TEXT,
	"created" DATETIME,
	PRIMARY KEY("id" AUTOINCREMENT)
	)`

var reviewTable = "review"

var reviewSchema = `
	CREATE TABLE "review"
        (
        "id" INTEGER,
        "card_id" INTEGER,
        "grade" INTEGER,
        "reviewed" DATETIME,
        PRIMARY KEY("id" AUTOINCREMENT),
        FOREIGN KEY (card_id) REFERENCES card(id)
	)`
