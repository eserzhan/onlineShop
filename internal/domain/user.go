package domain

import (
	"time"

)

// type User struct {
// 	ID           int   				  `json:"-" `
// 	Name         string               `json:"name" `
// 	Email        string               `json:"email" binding:"required"`
// 	Phone        string               `json:"phone" `
// 	Password     string               `json:"password" binding:"required"`
// 	Type 	     string 			  `json:"-"`
// 	RegisteredAt time.Time            `json:"-"`
// 	LastVisitAt  time.Time            `json:"lastVisitAt" bson:"lastVisitAt"`
// }

type User struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Phone        string    `json:"phone" db:"phone"`
	Password     string    `json:"password" db:"password_hash"`
	RegisteredAt time.Time `json:"registeredAt" db:"registered_at"`
	LastVisitAt  time.Time `json:"lastVisitAt" db:"last_visit_at"`
}