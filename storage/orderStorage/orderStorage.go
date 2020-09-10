package orderStorage

import (
	"errors"
	"fmt"
	"model/orderModel"
	"storage/shopDb"

	_ "github.com/lib/pq"
)

func GetOrdersForUser(userId string) []orderModel.Order {
	db := shopDb.ConnectToDb()
	defer db.Close()
	ordersSQL := "SELECT id, user_id, product_name FROM orders "
	ordersSQL += "WHERE user_id = $1;"
	rows, err := db.Query(ordersSQL, userId)
	var res []orderModel.Order
	if err != nil {
		fmt.Println(err.Error())
		return res
	}
	for rows.Next() {
		var order orderModel.Order
		err = rows.Scan(&order.Id, &order.UserId, &order.ProductName)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		res = append(res, order)
	}
	return res
}

func AddNewOrder(order orderModel.Order) (*orderModel.Order, error) {
	db := shopDb.ConnectToDb()
	defer db.Close()
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
	db := shopDb.ConnectToDb()
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

func UpdateOrder(id string, order orderModel.Order) (*orderModel.Order, error) {
	db := shopDb.ConnectToDb()
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
