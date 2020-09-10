package main

import (
	"controller/shopController"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"model/orderModel"
	"model/userModel"
	"net/http"
	"storage/shopDb"

	"github.com/gorilla/mux"
)

type ShopHamdler struct {
	handler http.Handler
}

func (l *ShopHamdler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" && r.URL.Path != "/register" {
		err := shopController.CheckAuthorization(w, r)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}
	}
	l.handler.ServeHTTP(w, r)
	// log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func NewShopHandler(handlerToWrap http.Handler) *ShopHamdler {
	return &ShopHamdler{handlerToWrap}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/login", login).Methods("POST")
	myRouter.HandleFunc("/register", register).Methods("POST")
	myRouter.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user_orders/{id}", userOrders)
	myRouter.HandleFunc("/user_order", addOrder).Methods("POST")
	myRouter.HandleFunc("/user_order/{id}", deleteOrder).Methods("DELETE")
	myRouter.HandleFunc("/user_order/{id}", updateOrder).Methods("PUT")
	wrappedMux := NewShopHandler(myRouter)
	log.Fatal(http.ListenAndServe(":8000", wrappedMux))
}

/*{
    "Id": "",
    "Name": "Adam",
    "Password": "12345"
}*/
func login(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user userModel.User
	json.Unmarshal(reqBody, &user)
	loginResult := shopController.Login(user)
	json.NewEncoder(w).Encode(loginResult)
}

func register(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user userModel.User
	json.Unmarshal(reqBody, &user)
	registerResult := shopController.Register(user)
	json.NewEncoder(w).Encode(registerResult)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	var user userModel.User
	json.Unmarshal(reqBody, &user)
	res := shopController.UpdateUser(id, user)
	json.NewEncoder(w).Encode(res)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	res := shopController.DeleteUser(id)
	json.NewEncoder(w).Encode(res)
}

func userOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	res := shopController.UserOrders(key)
	json.NewEncoder(w).Encode(res)
}

/*{
    "Id": "",
    "UserId": "1",
    "ProductName": "New product"
}*/
func addOrder(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var order orderModel.Order
	json.Unmarshal(reqBody, &order)
	res := shopController.AddOrder(order)
	json.NewEncoder(w).Encode(res)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	res := shopController.DeleteOrder(id)
	json.NewEncoder(w).Encode(res)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	var order orderModel.Order
	json.Unmarshal(reqBody, &order)
	res := shopController.UpdateOrder(id, order)
	json.NewEncoder(w).Encode(res)
}

func main() {
	shopDb.Init()
	handleRequests()
}
