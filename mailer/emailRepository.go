package mailer

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailRepo interface {
	Create(email *Email) error
	UpdateEmailStatus(primitive.ObjectID, Status) error
}
