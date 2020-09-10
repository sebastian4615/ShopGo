package userModel

type User struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
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

type DeleteResponse struct {
	Err  string `json:"Err"`
	IsOk bool   `json:"IsOk"`
}
