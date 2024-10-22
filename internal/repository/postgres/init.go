package postgres

import (
	"database/sql"
	"fmt"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Init(config *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		config.Database.User, config.Database.Password, config.Database.Name, config.Database.Port)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
