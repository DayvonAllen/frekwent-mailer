package mailer

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"myapp/database"
	"time"
)

type EmailRepoImpl struct {
	email Email
}

func (e EmailRepoImpl) Create(email *Email) error {
	conn := database.ConnectToDB()

	email.Template = "test"

	_, err := conn.EmailCollection.InsertOne(context.TODO(), email)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	fmt.Println("Sending: %v", email)

	go SendMessage(email)

	return nil
}

func (e EmailRepoImpl) UpdateEmailStatus(id primitive.ObjectID, status Status) error {
	conn := database.ConnectToDB()

	_, err := conn.EmailCollection.UpdateByID(context.TODO(), id, bson.D{{"$set",
		bson.D{{"updatedAt", time.Now()}, {"status", status}}}})

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("error processing data: %v", err)
	}

	return nil
}
