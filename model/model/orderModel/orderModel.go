package orderModel

type Order struct {
	Id          string `json:"Id"`
	UserId      string `json:"UserId"`
	ProductName string `json:"ProductName"`
}

type AddEditOrderResponse struct {
	Err      string `json:"Err"`
	NewOrder *Order `json:"NewOrder"`
}

type OrderListResponse struct {
	Err       string  `json:"Err"`
	OrderList []Order `json:"OrderList"`
}

type DeleteOrderResponse struct {
	Err  string `json:"Err"`
	IsOk bool   `json:"IsOk"`
}
