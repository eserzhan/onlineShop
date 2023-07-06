package domain

type CartItems struct {
	ID         int `json:"id", db:"id`
	Cart_id    int `json:"cart_id" db:"cart_id"`
	Product_id int `json:"product_id" db:"product_id"`
	Quantity   int `json:"quantity" db:"quantity"`
}
