package entity

import (
	"time"
)

type User struct {
	ID           string    `bson:"_id"`
	UserTypeID   Ref       `bson:"userTypeID"`
	Password     string    `bson:"password"`
	Username     string    `bson:"username"`
	CreatedAt    time.Time `bson:"createdAt"`
	RefreshToken string    `bson:"refreshToken"`
	Status       int       `bson:"status"`
}
