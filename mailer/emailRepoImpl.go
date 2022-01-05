package mailer

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"myapp/database"
	"time"
)

type EmailRepoImpl struct {
	email  Email
	emails []Email
}

func (e EmailRepoImpl) Create(email *Email) error {
	conn := database.Sess

	email.Template = "test"

	err := conn.DB(database.DB).C(database.EMAILS).Insert(&email)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	go SendMessage(email)

	return nil
}

func (e EmailRepoImpl) UpdateEmailStatus(id primitive.ObjectID, status Status) error {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.EMAILS).UpdateId(id, bson.M{"updatedAt": time.Now(), "status": status})

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("error processing data: %v", err)
	}

	return nil
}

func (e EmailRepoImpl) FindAllByStatus(status *Status) (*[]Email, error) {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.EMAILS).Find(bson.M{"status": status}).All(&e.emails)

	if err != nil {
		return nil, err
	}

	return &e.emails, nil
}
