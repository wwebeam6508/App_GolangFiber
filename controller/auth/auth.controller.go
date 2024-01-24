package controller

import (
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/auth"
	auth "PBD_backend_go/service/auth"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func LoginController(c *fiber.Ctx) error {
	var body model.LoginRequest
	err := c.BodyParser(&body)
	if err != nil {
		exception.PanicLogging(err)
	}

	// call the LoginService function
	result := auth.LoginService(body)

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
		exception.PanicLogging(err)
	}

	// call the RefreshTokenService function
	result := auth.RefreshTokenService(strings.Split(body.RefreshToken, " ")[1])
	tokenInpu := model.TokenInput{
		Token:  body.RefreshToken,
		UserID: result.UserID,
	}
	//call update refresh token
	auth.UpdateRefreshTokenService(tokenInpu)

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    bson.M{"accessToken": result.AccessToken, "userID": result.UserID},
	})
}

func FetchUserController(c *fiber.Ctx) error {
	userID := c.Params("userID")
	// call the FetchUserService function
	result := auth.FetchUserDataService(userID)

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}
