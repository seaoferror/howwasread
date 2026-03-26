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

	r := Repository{db: db}

	return &r
}
