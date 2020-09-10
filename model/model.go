package model

type User struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
}

type Order struct {
	Id          string `json:"Id"`
	UserId      string `json:"UserId"`
	ProductName string `json:"ProductName"`
}

type LoginResponse struct {
	Err   string `json:"Err"`
	Token string `json:"Token"`
}

type RegisterUpdateResponse struct {
	Err     string `json:"Err"`
	NewUser *User  `json:"NewUser"`
	Token   string `json:"Token"`
}

type AddEditOrderResponse struct {
	Err      string `json:"Err"`
	NewOrder *Order `json:"NewOrder"`
}
