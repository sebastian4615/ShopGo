package userStorage

import (
	"database/sql"
	"errors"
	"fmt"
	"model/userModel"
	"storage/shopDb"

	_ "github.com/lib/pq"
)

func CheckUserCredentials(name, password string) bool {
	user, err := GetUserByName(name)
	if err != nil {
		return false
	}
	return user.Password == password
}

func AddNewUser(name, password string) (*userModel.User, error) {
	if password == "" || name == "" {
		return nil, errors.New("Password and user name cannot be empty!")
	}
	db := shopDb.ConnectToDb()
	defer db.Close()
	_, notExist := getUserByName(db, name)
	if notExist == nil {
		return nil, errors.New("User alredy exist")
	}
	id, err := addNewUser(db, name, password)
	if err != nil {
		return nil, err
	}
	return &userModel.User{Id: id, Name: name, Password: password}, nil
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

func GetUserByName(name string) (*userModel.User, error) {
	db := shopDb.ConnectToDb()
	defer db.Close()
	return getUserByName(db, name)
}

func getUserByName(db *sql.DB, name string) (*userModel.User, error) {
	userSql := "SELECT id, name, password FROM users WHERE name = $1"
	var user userModel.User
	err := db.QueryRow(userSql, name).Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserById(id string) (*userModel.User, error) {
	db := shopDb.ConnectToDb()
	defer db.Close()
	return getUserById(db, id)
}

func getUserById(db *sql.DB, id string) (*userModel.User, error) {
	userSql := "SELECT id, name, password FROM users WHERE id = $1"
	var user userModel.User
	err := db.QueryRow(userSql, id).Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id string, user userModel.User) (*userModel.User, error) {
	if user.Password == "" || user.Name == "" {
		return nil, errors.New("Password and user name cannot be empty!")
	}
	db := shopDb.ConnectToDb()
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
	db := shopDb.ConnectToDb()
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
