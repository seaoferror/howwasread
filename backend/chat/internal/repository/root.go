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
		id BINARY(16) NOT NULL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS payload (
		id BINARY(16) NOT NULL PRIMARY KEY,
		to_id BINARY(16) NOT NULL,
	    from_id BINARY(16) NOT NULL,
	    content_type VARCHAR(50) NOT NULL,
		content TEXT NOT NULL,
		INDEX idx_message_to_id_id (to_id, id)
	);`)
	if err != nil {
		panic(err)
	}
	//json fromId, type, content, created_at
	// sender name should be updated whenever member change their profile, so no need to send them and save in sqlite

	r := Repository{db: db}

	return &r
}
