package domain 

type Cart struct {
	ID int `json:"id" db:"id"`
	UserId string `json:"user_id" db:"user_id"`
}