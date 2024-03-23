package model

type OrderItemDTO struct {
	DishName string `json:"dishName"`
	Quantity int    `json:"quantity"`
	Comment  string `json:"comment"`
}

type ClientOrderDTO struct {
	ClientId string         `json:"clientId"`
	Items    []OrderItemDTO `json:"items"`
}

type ClientOrderResponseDTO struct {
	Comment string `json:"comment"`
}
