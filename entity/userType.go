package entity

import (
	"reflect"
	"time"
)

type PermissionDetail struct {
	CanEdit   bool `json:"canEdit" bson:"canEdit"`
	CanRemove bool `json:"canRemove" bson:"canRemove"`
	CanView   bool `json:"canView" bson:"canView"`
}

type Permissions struct {
	User          PermissionDetail `json:"user" bson:"user"`
	UserType      PermissionDetail `json:"userType" bson:"userType"`
	Project       PermissionDetail `json:"project" bson:"project"`
	Expense       PermissionDetail `json:"expense" bson:"expense"`
	Customer      PermissionDetail `json:"customer" bson:"customer"`
	Employee      PermissionDetail `json:"employee" bson:"employee"`
	Location      PermissionDetail `json:"location" bson:"location"`
	InventoryType PermissionDetail `json:"inventoryType" bson:"inventoryType"`
	Inventory     PermissionDetail `json:"inventory" bson:"inventory"`
	Wage          PermissionDetail `json:"wage" bson:"wage"`
	Requisition   PermissionDetail `json:"requisition" bson:"requisition"`
}

type UserType struct {
	ID         string      `json:"_id" bson:"_id"`
	Name       string      `json:"name" bson:"name"`
	Permission Permissions `json:"permission" bson:"permission"`
	Status     int         `json:"status" bson:"status"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
}

func NewPermission() Permissions {
	permission := Permissions{}

	vPermission := reflect.ValueOf(&permission).Elem()

	for i := 0; i < vPermission.NumField(); i++ {
		vPermission.Field(i).Set(reflect.ValueOf(PermissionDetail{
			CanEdit:   false,
			CanRemove: false,
			CanView:   false,
		}))
	}
	return permission
}
