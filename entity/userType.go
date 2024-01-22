package entity

import (
	"time"
)

type PermissionDetail struct {
	CanEdit   bool `bson:"canEdit"`
	CanRemove bool `bson:"canRemove"`
	CanView   bool `bson:"canView"`
}

type Permissions struct {
	User     PermissionDetail `bson:"user"`
	UserType PermissionDetail `bson:"userType"`
	Project  PermissionDetail `bson:"project"`
	Expense  PermissionDetail `bson:"expense"`
}

type UserType struct {
	ID         string      `bson:"_id"`
	Name       string      `bson:"name"`
	Permission Permissions `bson:"permission"`
	Status     int         `bson:"status"`
	CreatedAt  time.Time   `bson:"createdAt"`
}
