package beurse

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"Email" bson:"email"`
	Password string `json:"password" bson:"password"`
	PhoneNumber string             `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	// Role     string `json:"role,omitempty" bson:"role,omitempty"`
}

type Device struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Topic string             `json:"topic" bson:"topic"`
	User  string             `json:"user" bson:"user"`
	Status bool			     `json:"status" bson:"status"`
}

type History struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Topic     string             `json:"topic" bson:"topic"`
	Payload   string             `json:"payload" bson:"payload"`
	User      string             `json:"user" bson:"user"`
	CreatedAt time.Time			 `json:"created_at" bson:"created_at"`
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

type HistoryResponse struct {
    Status  bool      `json:"status"`
    Message string    `json:"message"`
    Data    []History `json:"data"`
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type Otp struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	OTP       string             `json:"otp" bson:"otp"`
	ExpiredAt int64              `json:"expiredat" bson:"expiredat"`
	Status    bool               `json:"status" bson:"status"`
}

type ResetPassword struct {
	Email    string `json:"email" bson:"email"`
	OTP      string `json:"otp" bson:"otp"`
	Password string `json:"password" bson:"password"`
}