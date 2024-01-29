package model

import (
	"PBD_backend_go/entity"
	authmodel "PBD_backend_go/model/auth"
)

type PermissionResult struct {
	UserType authmodel.UserTypeResult `json:"userType" bson:"userType"`
}

type PermissionUserType struct {
	Permission entity.Permissions `json:"permission" bson:"permission"`
}

type PermissionInput struct {
	GroupName string
	Name      string
}
