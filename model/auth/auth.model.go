package model

import (
	"PBD_backend_go/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRequest struct {
	Username string
	Password string
}

type RefreshTokenRequest struct {
	RefreshToken string
}

type UserIDInput struct {
	UserID string
}

type TokenInput struct {
	Token  string
	UserID string
}

type UserTypeResult struct {
	UserTypeID primitive.ObjectID `json:"userTypeID" bson:"userTypeID"`
	Name       *string            `json:"name" bson:"name"`
	Permission entity.Permissions `json:"permission" bson:"permission"`
	Rank       *int32             `json:"rank" bson:"rank"`
}

type UserProfileResult struct {
	UserID   primitive.ObjectID `json:"userID" bson:"userID"`
	Username *string            `json:"username" bson:"username"`
	Password *string            `json:"password" bson:"password"`
	UserType UserTypeResult     `json:"userType" bson:"userType"`
}

type FetchUserResult struct {
	UserData      UserProfileResult  `json:"userData" bson:"userData"`
	PrePermission entity.Permissions `json:"prePermission" bson:"prePermission"`
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
}

type ChangePasswordInput struct {
	UserID          string `json:"userID" bson:"userID"`
	Password        string `json:"Password" bson:"Password"`
	ConfirmPassword string `json:"ConfirmPassword" bson:"ConfirmPassword"`
}
