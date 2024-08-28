package model

type MenuDto struct {
	Items []*MenuItemDto `json:"items"`
}

type MenuItemDto struct {
	Name        string   `json:"name"`
	Price       float32  `json:"price"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}
