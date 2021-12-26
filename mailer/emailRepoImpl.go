package mailer

import (
	"context"
	"errors"
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

func (e EmailRepoImpl) FindAllByStatus(status *Status) (*[]Email, error) {
	conn := database.ConnectToDB()

	cur, err := conn.EmailCollection.Find(context.TODO(), bson.D{{"status", status}})

	if err != nil {
		return nil, errors.New("error finding email")
	}

	if err = cur.All(context.TODO(), &e.emails); err != nil {
		panic(err)
	}

	if e.emails == nil {
		return nil, errors.New("no emails found with that status")
	}

	return &e.emails, nil
}
