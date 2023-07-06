package domain

type Order struct {
	ID int `json:"id" db:"id"`
	User_id int `json:"user_id" db:"user_id"`
	Status string `json:"status" db:"status"`
	Delivery_method string `json:"delivery_method" db:"delivery_method"`
	Payment_method string `json:"payment_method" db:"payment_method"`
}