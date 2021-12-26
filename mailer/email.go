package mailer

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Status int32

const (
	Success Status = iota
	Pending
	Failed
)

type Email struct {
	Id            primitive.ObjectID `bson:"_id" json:"id"`
	Type          string             `bson:"type" json:"type"`
	CustomerEmail string             `bson:"customerEmail" json:"customerEmail"`
	From          string             `bson:"from" json:"from"`
	Content       string             `bson:"content" json:"content"`
	Subject       string             `bson:"subject" json:"subject"`
	Status        Status             `bson:"status" json:"status"`
	CreatedAt     time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}