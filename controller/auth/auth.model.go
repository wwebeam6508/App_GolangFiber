package controller

type LoginRequest struct {
	username string `json:"username"`
	password string `json:"password"`
}

