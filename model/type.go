package beurse

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"Email" bson:"email"`
	// Role     string `json:"role,omitempty" bson:"role,omitempty"`
}
type Device struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Topic string             `json:"topic" bson:"topic"`
	User string             `json:"user" bson:"user"`
}

type Credential struct {
	Status  bool   `json:"status" bson:"status"`
	Token   string `json:"token,omitempty" bson:"token,omitempty"`
	Message string `json:"message,omitempty" bson:"message,omitempty"`
}

type DeviceResponse struct {
    Status  bool     `json:"status"`
    Message string   `json:"message"`
    Data    []Device `json:"data"`
}