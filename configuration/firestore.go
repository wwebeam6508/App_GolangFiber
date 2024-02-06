package configuration

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func ConnectToFirestore() (*firestore.Client, error) {
	// Set up the Firebase credentials file path
	opt := option.WithCredentialsFile("firebase_credentials.json")

	// Initialize the Firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
		return nil, err
	}

	// Get a Firestore client
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return nil, err
	}

	return client, nil
}

func ConnectToStorage() (*storage.Client, error) {
	// Set up the Firebase credentials file path
	opt := option.WithCredentialsFile("firebase_credentials.json")

	client, err := storage.NewClient(context.Background(), opt)
	if err != nil {
		log.Fatalf("Failed to create Storage client: %v", err)
		return nil, err
	}

	return client, nil

}
