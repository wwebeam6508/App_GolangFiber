package entity

import (
	"time"
)

type PermissionDetail struct {
	CanEdit   bool `json:"canEdit" bson:"canEdit"`
	CanRemove bool `json:"canRemove" bson:"canRemove"`
	CanView   bool `json:"canView" bson:"canView"`
}

type Permissions struct {
	User     PermissionDetail `json:"user" bson:"user"`
	UserType PermissionDetail `json:"userType" bson:"userType"`
	Project  PermissionDetail `json:"project" bson:"project"`
	Expense  PermissionDetail `json:"expense" bson:"expense"`
}

type UserType struct {
	ID         string      `json:"_id" bson:"_id"`
	Name       string      `json:"name" bson:"name"`
	Permission Permissions `json:"permission" bson:"permission"`
	Status     int         `json:"status" bson:"status"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
}
