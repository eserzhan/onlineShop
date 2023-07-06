package domain

type Product struct {
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Price float32 `json:"price" db:"price"`
	Quantity int `json:"quantity" db:"quantity"`
}

type UpdateProduct struct {
	Name *string `json:"name"`
	Description *string `json:"description"`
	Price *float32 `json:"price"`
	Quantity *int `json:"quantity"`
}
