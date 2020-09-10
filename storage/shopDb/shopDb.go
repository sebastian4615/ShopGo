package shopDb

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "ShopDB"
)

func Init() {
	db := ConnectToDb()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	createTables(db)
}

func ConnectToDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func createTables(db *sql.DB) {
	query := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"; "
	query += "CREATE TABLE IF NOT EXISTS users ("
	query += "id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), "
	query += "name VARCHAR, "
	query += "password VARCHAR "
	query += ");"

	query += "CREATE TABLE IF NOT EXISTS orders ("
	query += "id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), "
	query += "user_id uuid, "
	query += "product_name VARCHAR "
	query += ");"

	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}
