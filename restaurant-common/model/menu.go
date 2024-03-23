package model

type Menu struct {
	Items []MenuItem
}

type MenuItem struct {
	Name        string   `json:"name"`
	Price       float32  `json:"price"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}
