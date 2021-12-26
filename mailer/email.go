package mailer

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Status int32
type EmailType int32

const (
	Success Status = iota
	Pending
	Failed
)

const (
	CouponPromotion EmailType = iota
	Purchased
	CustomerInteraction
)

type Email struct {
	Id            primitive.ObjectID `bson:"_id" json:"id"`
	Type          EmailType          `bson:"type" json:"type"`
	CustomerEmail string             `bson:"customerEmail" json:"customerEmail"`
	From          string             `bson:"from" json:"from"`
	Content       string             `bson:"content" json:"content"`
	Subject       string             `bson:"subject" json:"subject"`
	Status        Status             `bson:"status" json:"status"`
	Template      string             `bson:"template" json:"template"`
	Attachments   []string           `bson:"attachments" json:"attachments"`
	Data          interface{}        `bson:"data" json:"data"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}
