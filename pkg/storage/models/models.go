package models

import "time"

// User представляет модель пользователя
type User struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Inventory представляет модель инвентаря пользователя
type Inventory struct {
	ID     string  `json:"id"`
	UserID string  `json:"userid"`
	Coins  int     `json:"coins"`
	Merch  []Merch `json:"merch"`
}

// Merch представляет модель товара
type Merch struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// Transaction представляет модель транзакции
type Transaction struct {
	ID        int       `json:"id"`
	FromUser  string    `json:"from_user"`
	ToUser    string    `json:"to_user"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type BuyItem struct {
	UserID string
	Item   string
}

type InventoryOut struct {
	Name     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type FromUser struct {
	FromUser  string    `json:"fromUser"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type ToUser struct {
	ToUser    string    `json:"toUser"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type CoinHistory struct {
	Received []FromUser `json:"received"`
	Sent     []ToUser   `json:"sent"`
}

type Info struct {
	Coins       int            `json:"coins"`
	Inventory   []InventoryOut `json:"inventory"`
	CoinHistory CoinHistory    `json:"coinHistory"`
}
