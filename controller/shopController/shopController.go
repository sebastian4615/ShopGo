package shopController

import (
	"errors"
	"fmt"
	"model/orderModel"
	"model/userModel"
	"net/http"
	"storage/orderStorage"
	"storage/userStorage"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CheckAuthorization(w http.ResponseWriter, r *http.Request) error {
	if r.Header["Token"] != nil {
		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			claims := token.Claims.(jwt.MapClaims)
			userName, isString := claims["client"].(string)
			if !isString {
				return nil, errors.New("Authorization failed")
			}
			user, userError := userStorage.GetUserByName(userName)
			if userError != nil {
				return nil, errors.New("Authorization failed")
			}
			return []byte(user.Password), nil
		})

		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		if token.Valid {
			return nil
		}
		return errors.New("Authorization failed")
	} else {
		return errors.New("Wrong header format")
	}
}

func GenerateJWT(userName, password string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = userName
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	var mySigningKey = []byte(password)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func Login(user userModel.User) userModel.LoginResponse {
	if userStorage.CheckUserCredentials(user.Name, user.Password) {
		token, err := GenerateJWT(user.Name, user.Password)
		if err != nil {
			// json.NewEncoder(w).Encode(userModel.LoginResponse{Err: err.Error(), Token: ""})
			return userModel.LoginResponse{Err: err.Error(), Token: ""}
		}
		return userModel.LoginResponse{Err: "", Token: token}
	}
	return userModel.LoginResponse{Err: "Wrong credentials", Token: ""}
}

func Register(user userModel.User) userModel.RegisterUpdateResponse {
	newUser, err := userStorage.AddNewUser(user.Name, user.Password)
	return createUserResponse(newUser, err)
}

func UpdateUser(id string, user userModel.User) userModel.RegisterUpdateResponse {
	newUser, err := userStorage.UpdateUser(id, user)
	return createUserResponse(newUser, err)
}

func DeleteUser(id string) userModel.DeleteResponse {
	res := userStorage.DeleteUser(id)
	if res {
		return userModel.DeleteResponse{Err: "", IsOk: true}
	}
	return userModel.DeleteResponse{Err: "Delete user failed", IsOk: false}
}

func createUserResponse(newUser *userModel.User, err error) userModel.RegisterUpdateResponse {
	if err != nil {
		return userModel.RegisterUpdateResponse{Err: err.Error(), NewUser: nil, Token: ""}
	} else {
		token, err := GenerateJWT(newUser.Name, newUser.Password)
		if err != nil {
			return userModel.RegisterUpdateResponse{Err: err.Error(), NewUser: nil, Token: ""}
		}
		return userModel.RegisterUpdateResponse{Err: "", NewUser: newUser, Token: token}
	}
}

func UserOrders(userId string) []orderModel.Order {
	return orderStorage.GetOrdersForUser(userId)
}

func AddOrder(order orderModel.Order) orderModel.AddEditOrderResponse {
	_, err := userStorage.GetUserById(order.UserId)
	if err != nil {
		return orderModel.AddEditOrderResponse{Err: "Wrong user ID", NewOrder: nil}
	}
	newOrder, err := orderStorage.AddNewOrder(order)
	if err != nil {
		return orderModel.AddEditOrderResponse{Err: err.Error(), NewOrder: nil}
	}
	return orderModel.AddEditOrderResponse{Err: "", NewOrder: newOrder}
}

func DeleteOrder(id string) orderModel.DeleteOrderResponse {
	res := orderStorage.DeleteOrder(id)
	if res {
		return orderModel.DeleteOrderResponse{Err: "", IsOk: true}
	}
	return orderModel.DeleteOrderResponse{Err: "Delete order fail", IsOk: false}

}
func UpdateOrder(id string, order orderModel.Order) orderModel.AddEditOrderResponse {
	newOrder, err := orderStorage.UpdateOrder(id, order)
	if err != nil {
		return orderModel.AddEditOrderResponse{Err: err.Error(), NewOrder: nil}
	}
	return orderModel.AddEditOrderResponse{Err: "", NewOrder: newOrder}
}
