package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"model"

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
	db := connectToDb()
	defer db.Close()
	err := db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	createTables(db)
}

func connectToDb() *sql.DB {
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
		fmt.Println(err.Error())
	}
}

func CheckUserCredentials(name, password string) bool {
	user, err := GetUserByName(name)
	if err != nil {
		return false
	}
	return user.Password == password
}

func AddNewUser(name, password string) (*model.User, error) {
	if password == "" || name == "" {
		return nil, errors.New("Password and user name cannot be empty!")
	}
	db := connectToDb()
	defer db.Close()
	_, notExist := getUserByName(db, name)
	if notExist == nil {
		return nil, errors.New("User alredy exist")
	}
	id, err := addNewUser(db, name, password)
	if err != nil {
		return nil, err
	}
	return &model.User{Id: id, Name: name, Password: password}, nil
}

func addNewUser(db *sql.DB, name, password string) (string, error) {
	userSQL := "INSERT INTO users (name, password) VALUES($1, $2) RETURNING id;"
	stmt, err := db.Prepare(userSQL)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var userID string
	err = stmt.QueryRow(name, password).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func GetUserByName(name string) (*model.User, error) {
	db := connectToDb()
	defer db.Close()
	return getUserByName(db, name)
}

func getUserByName(db *sql.DB, name string) (*model.User, error) {
	userSql := "SELECT id, name, password FROM users WHERE name = $1"
	var user model.User
	err := db.QueryRow(userSql, name).Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func getUserById(db *sql.DB, id string) (*model.User, error) {
	userSql := "SELECT id, name, password FROM users WHERE id = $1"
	var user model.User
	err := db.QueryRow(userSql, id).Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id string, user model.User) (*model.User, error) {
	if user.Password == "" || user.Name == "" {
		return nil, errors.New("Password and user name cannot be empty!")
	}
	db := connectToDb()
	defer db.Close()
	userSQL := "UPDATE users "
	userSQL += "SET name = $2, password = $3 "
	userSQL += "WHERE id = $1;"
	res, err := db.Exec(userSQL, id, user.Name, user.Password)
	if err != nil {
		return nil, errors.New("Update user fail")
	}
	rowsNum, err := res.RowsAffected()
	if err != nil && rowsNum == 0 {
		return nil, errors.New("Update user fail")
	}
	user.Id = id
	return &user, nil
}

func DeleteUser(id string) bool {
	db := connectToDb()
	defer db.Close()
	userSQL := "DELETE FROM users WHERE id = $1;"
	res, err := db.Exec(userSQL, id)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rowsNum, err := res.RowsAffected()
	if err == nil && rowsNum > 0 {
		ordersSQL := "DELETE FROM orders WHERE user_id = $1;"
		_, err := db.Exec(ordersSQL, id)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		return true
	}
	return false
}

func GetOrdersForUser(userId string) []model.Order {
	db := connectToDb()
	defer db.Close()
	ordersSQL := "SELECT id, user_id, product_name FROM orders "
	ordersSQL += "WHERE user_id = $1;"
	rows, err := db.Query(ordersSQL, userId)
	var res []model.Order
	if err != nil {
		fmt.Println(err.Error())
		return res
	}
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.Id, &order.UserId, &order.ProductName)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		res = append(res, order)
	}
	return res
}

func AddNewOrder(order model.Order) (*model.Order, error) {
	db := connectToDb()
	defer db.Close()
	_, err := getUserById(db, order.UserId)
	if err != nil {
		return nil, errors.New("Wrong user ID")
	}
	orderSQL := "INSERT INTO orders (user_id, product_name) VALUES($1, $2) RETURNING id;"
	stmt, err := db.Prepare(orderSQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var orderID string
	err = stmt.QueryRow(order.UserId, order.ProductName).Scan(&orderID)
	if err != nil {
		return nil, err
	}
	order.Id = orderID
	return &order, nil
}

func DeleteOrder(id string) bool {
	db := connectToDb()
	defer db.Close()
	ordersSQL := "DELETE FROM orders WHERE id = $1;"
	res, err := db.Exec(ordersSQL, id)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rowsNum, err := res.RowsAffected()
	return err == nil && rowsNum > 0
}

func UpdateOrder(id string, order model.Order) (*model.Order, error) {
	db := connectToDb()
	defer db.Close()
	ordersSQL := "UPDATE orders "
	ordersSQL += "SET product_name = $3 "
	ordersSQL += "WHERE id = $1 AND user_id = $2;"
	res, err := db.Exec(ordersSQL, id, order.UserId, order.ProductName)
	if err != nil {
		return nil, errors.New("Update order fail")
	}
	rowsNum, err := res.RowsAffected()
	if err != nil && rowsNum == 0 {
		return nil, errors.New("Update order fail")
	}
	return &order, nil
}
