package controller

import (
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/auth"
	authservice "PBD_backend_go/service/auth"

	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func LoginController(c *fiber.Ctx) error {
	var body model.LoginRequest
	err := c.BodyParser(&body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	// call the LoginService function
	result, err := authservice.LoginService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	//update refresh token
	tokenInput := model.TokenInput{
		Token:  result.RefreshToken,
		UserID: result.UserProfile.UserID.Hex(),
	}
	err = authservice.UpdateRefreshTokenService(tokenInput)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func RefreshTokenController(c *fiber.Ctx) error {
	var body model.RefreshTokenRequest
	err := c.BodyParser(&body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	// call the RefreshTokenService function
	result, err := authservice.RefreshTokenService(strings.Split(body.RefreshToken, " ")[1])
	tokenInpu := model.TokenInput{
		Token:  body.RefreshToken,
		UserID: result.UserID,
	}
	if err != nil {
		return exception.ErrorHandler(c, exception.AccessDenialError{Message: err.Error()})
	}
	//call update refresh token
	err = authservice.UpdateRefreshTokenService(tokenInpu)
	if err != nil {
		return exception.ErrorHandler(c, exception.AccessDenialError{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    bson.M{"accessToken": result.AccessToken, "userID": result.UserID},
	})
}

func FetchUserController(c *fiber.Ctx) error {
	userID := c.Query("userID")
	// call the FetchUserService function
	result, err := authservice.FetchUserDataService(userID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func LogoutController(c *fiber.Ctx) error {
	var body model.UserIDInput
	err := c.BodyParser(&body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	err = authservice.RemoveRefreshTokenService(body.UserID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}

func ChangePasswordController(c *fiber.Ctx) error {
	var body model.ChangePasswordRequest
	err := c.BodyParser(&body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	//verify JWT by headers.Authorization and split to get token
	split := strings.Split(c.Get("Authorization"), " ")
	//check if split is 2 or not
	if len(split) != 2 {
		return exception.ErrorHandler(c, exception.ValidationError{Message: "invalid token"})
	}
	token := split[1]
	claims, err := authservice.VerifyJWT(token)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	//get userID from claims
	userID := claims.Claims.(jwt.MapClaims)["data"].(map[string]interface{})["userID"].(string)
	//create input for ChangePasswordService
	input := model.ChangePasswordInput{
		UserID:          userID,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	}
	// call the ChangePasswordService function
	authservice.ChangePasswordDataService(input)

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}
