package models

import "time"

type Menu struct {
	ID              string      `bson:"_id" json:"_id"`
	Name            string      `bson:"name" json:"name"`
	RestaurantID    string      `bson:"restaurant_id" json:"restaurant_id"`
	Products        interface{} `bson:"products" json:"products"`
	Attributes      interface{} `bson:"attributes" json:"attributes"`
	AttributesGroup interface{} `bson:"attributes_groups" json:"attributes_groups"`
	CreatedAt       time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `bson:"updated_at" json:"updated_at"`
}
