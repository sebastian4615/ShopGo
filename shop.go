package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"model"
	"net/http"
	"storage"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Shop HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.Handle("/user_orders/{id}", isAuthorized(userOrders))
	myRouter.Handle("/user_order", isAuthorized(addOrder)).Methods("POST")
	myRouter.Handle("/user_order/{id}", isAuthorized(deleteOrder)).Methods("DELETE")
	myRouter.Handle("/user_order/{id}", isAuthorized(updateOrder)).Methods("PUT")
	myRouter.Handle("/user/{id}", isAuthorized(updateUser)).Methods("PUT")
	myRouter.Handle("/user/{id}", isAuthorized(deleteUser)).Methods("DELETE")
	myRouter.HandleFunc("/login", login).Methods("POST")
	myRouter.HandleFunc("/register", register).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}

/*{
    "Id": "",
    "Name": "Adam",
    "Password": "12345"
}*/
func login(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user model.User
	json.Unmarshal(reqBody, &user)
	if storage.CheckUserCredentials(user.Name, user.Password) {
		token, err := generateJWT(user.Name, user.Password)
		if err != nil {
			json.NewEncoder(w).Encode(model.LoginResponse{Err: err.Error(), Token: ""})
			return
		}
		json.NewEncoder(w).Encode(model.LoginResponse{Err: "", Token: token})
		return
	}
	json.NewEncoder(w).Encode(model.LoginResponse{Err: "Wrong credentials", Token: ""})
}

func register(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user model.User
	json.Unmarshal(reqBody, &user)
	newUser, err := storage.AddNewUser(user.Name, user.Password)
	sendUser(w, newUser, err)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user model.User
	json.Unmarshal(reqBody, &user)
	newUser, err := storage.UpdateUser(id, user)
	sendUser(w, newUser, err)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	res := storage.DeleteUser(id)
	if res {
		json.NewEncoder(w).Encode("User deleted")
	} else {
		json.NewEncoder(w).Encode("User order fail")
	}
}

func sendUser(w http.ResponseWriter, newUser *model.User, err error) {
	if err != nil {
		json.NewEncoder(w).Encode(model.RegisterUpdateResponse{Err: err.Error(), NewUser: nil, Token: ""})
	} else {
		token, err := generateJWT(newUser.Name, newUser.Password)
		if err != nil {
			json.NewEncoder(w).Encode(model.RegisterUpdateResponse{Err: err.Error(), NewUser: nil, Token: ""})
			return
		}
		json.NewEncoder(w).Encode(model.RegisterUpdateResponse{Err: "", NewUser: newUser, Token: token})
	}
}

func userOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	orders := storage.GetOrdersForUser(key)
	json.NewEncoder(w).Encode(orders)
}

/*{
    "Id": "",
    "UserId": "1",
    "ProductName": "New product"
}*/
func addOrder(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var order model.Order
	json.Unmarshal(reqBody, &order)
	newOrder, err := storage.AddNewOrder(order)
	if err != nil {
		json.NewEncoder(w).Encode(model.AddEditOrderResponse{Err: err.Error(), NewOrder: nil})
	} else {
		json.NewEncoder(w).Encode(model.AddEditOrderResponse{Err: "", NewOrder: newOrder})
	}
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	res := storage.DeleteOrder(id)
	if res {
		json.NewEncoder(w).Encode("Order deleted")
	} else {
		json.NewEncoder(w).Encode("Delete order fail")
	}
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	var order model.Order
	json.Unmarshal(reqBody, &order)
	newOrder, err := storage.UpdateOrder(id, order)
	if err != nil {
		json.NewEncoder(w).Encode(model.AddEditOrderResponse{Err: err.Error(), NewOrder: nil})
	} else {
		json.NewEncoder(w).Encode(model.AddEditOrderResponse{Err: "", NewOrder: newOrder})
	}
}

func generateJWT(userName, password string) (string, error) {
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

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				user, userError := storage.GetUserByName(userName)
				if userError != nil {
					return nil, errors.New("Authorization failed")
				}
				return []byte(user.Password), nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func main() {
	storage.Init()
	handleRequests()
}
