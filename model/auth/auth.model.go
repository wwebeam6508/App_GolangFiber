package model

import (
	"PBD_backend_go/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRequest struct {
	Username string
	Password string
}

type TokenRequest struct {
	Token  string
	UserID string
}

type UserTypeResult struct {
	UserTypeID  primitive.ObjectID `json:"userTypeID" bson:"userTypeID"`
	Name        *string            `json:"name" bson:"name"`
	Permissions entity.Permissions `json:"permissions" bson:"permissions"`
}

type UserProfileResult struct {
	UserID   primitive.ObjectID `json:"userID" bson:"userID"`
	Username *string            `json:"username" bson:"username"`
	Password *string            `json:"password" bson:"password"`
	UserType UserTypeResult     `json:"userType" bson:"userType"`
}

type UserResult struct {
	AccessToken  string            `json:"accessToken" bson:"accessToken"`
	RefreshToken string            `json:"refreshToken" bson:"refreshToken"`
	UserProfile  UserProfileResult `json:"userProfile" bson:"userProfile"`
}

type RefreshTokenResult struct {
	AccessToken string `json:"accessToken" bson:"accessToken"`
	UserID      string `json:"userID" bson:"userID"`
}

type ChangePasswordRequest struct {
	Password        string `json:"Password" bson:"Password"`
	ConfirmPassword string `json:"ConfirmPassword" bson:"ConfirmPassword"`
	UserID          string `json:"UserID" bson:"UserID"`
}
