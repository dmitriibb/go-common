package model

type ReadyOrderItem struct {
	OrderId  int    `json:"orderId"`
	ItemId   int    `json:"itemId"`
	DishName string `json:"dishName"`
	Comment  string `json:"comment"`
	Payload  string `json:"payload"`
}
