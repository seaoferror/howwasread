package repository

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	db, err := sql.Open("mysql", os.Getenv("MY_SQL_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS member (
		id CHAR(36) NOT NULL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS message (
		id CHAR(36) NOT NULL PRIMARY KEY,
		to_id CHAR(36) NOT NULL,
		content JSON NOT NULL,
		CONSTRAINT fk_message_to_member FOREIGN KEY (to_id) REFERENCES member(id),
		INDEX INDEX idx_message_to_id (to_id);
	);`)
	if err != nil {
		panic(err)
	}
	//json fromId, type, content, created_at
	// sender name should be updated whenever member change their profile, so no need to send them and save in sqlite

	r := Repository{db: db}

	return &r
}
