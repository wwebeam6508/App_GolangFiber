package middleware

import (
	"PBD_backend_go/exception"
	model "PBD_backend_go/model"
	service "PBD_backend_go/service"
	authservice "PBD_backend_go/service/auth"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func RankCheck(c *fiber.Ctx) error {
	//get body userID and userTypeID from body
	var body struct {
		UserID     string `json:"userID" bson:"userID"`
		UserTypeID string `json:"userTypeID" bson:"userTypeID"`
	}
	err := c.BodyParser(&body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	token := c.Get("Authorization")
	splitToken := strings.Split(token, " ")
	if len(splitToken) < 2 {
		return exception.ValidationError{Message: "invalid token"}
	}
	claims, err := authservice.VerifyJWT(splitToken[1])
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	userData := claims.Claims.(jwt.MapClaims)["data"].(map[string]interface{})
	//check if userID is empty even it spacing
	if body.UserID != "" {
		err := againistOther(userData, body.UserID)
		if err != nil {
			return exception.ErrorHandler(c, err)
		}
	} else if body.UserTypeID != "" {
		//get userTypeID from claim
		err := againistOtherType(userData, body.UserTypeID)
		if err != nil {
			return exception.ErrorHandler(c, err)
		}
	}
	return c.Next()
}

func againistOther(userData map[string]interface{}, userID string) error {
	selfUserID := userData["userID"].(string)
	//check if userID from body is equal to userID from claim
	if userID == selfUserID {
		return exception.ValidationError{Message: "cannot change your own data"}
	}
	//get userTypeID from input.SelfID
	rank, err := service.GetUserRankByUserIDService(model.GetUserTypeByUserIDInput{UserID: userID})
	if err != nil {
		return err
	}
	selfRank, err := service.GetUserRankByUserIDService(model.GetUserTypeByUserIDInput{UserID: selfUserID})
	if err != nil {
		return err
	}
	//check if userTypeID is super admin
	if selfRank.Rank <= rank.Rank {
		return exception.UnauthorizedError{Message: "cannot change rank higher than or equal to your rank"}
	}

	return nil
}

func againistOtherType(userData map[string]interface{}, userTypeID string) error {
	selfUserTypeID := userData["userType"].(map[string]interface{})["userTypeID"].(string)
	//get rank from userTypeID
	if userTypeID == selfUserTypeID {
		return exception.ValidationError{Message: "cannot change your own data"}
	}
	//get userTypeID from input.SelfID
	rank, err := service.GetUserRankByUserTypeIDService(model.GetUserRankByUserTypeIDInput{UserTypeID: userTypeID})
	if err != nil {
		return err
	}
	selfRank, err := service.GetUserRankByUserTypeIDService(model.GetUserRankByUserTypeIDInput{UserTypeID: selfUserTypeID})
	if err != nil {
		return err
	}
	//check if userTypeID is super admin
	if selfRank.Rank <= rank.Rank {
		return exception.UnauthorizedError{Message: "cannot change rank higher than or equal to your rank"}
	}

	return nil
}
