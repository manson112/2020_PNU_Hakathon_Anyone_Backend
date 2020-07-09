package fb

import (
	"base/models"
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client
var authClient *auth.Client
var messagingClient *messaging.Client

// InitFirebaseAdminSDK :: Setting up firebase admin functions
func InitFirebaseAdminSDK() {
	credential := models.FirebaseCredential{
		Type:                    os.Getenv("FB_ADMINSDK_TYPE"),
		ProjectID:               os.Getenv("FB_ADMINSDK_PROJECT_ID"),
		PrivateKeyID:            os.Getenv("FB_ADMINSDK_PRIVATE_KEY_ID"),
		PrivateKey:              os.Getenv("FB_ADMINSDK_PRIVATE_KEY"),
		ClientEmail:             os.Getenv("FB_ADMINSDK_CLIENT_EMAIL"),
		ClientID:                os.Getenv("FB_ADMINSDK_CLIENT_ID"),
		AuthURI:                 os.Getenv("FB_ADMINSDK_AUTH_URI"),
		TokenURI:                os.Getenv("FB_ADMINSDK_TOKEN_URI"),
		AuthProviderX509CertURL: os.Getenv("FB_ADMINSDK_AUTH_PROVIDER_X509_CERT_URL"),
		ClientX509CertURL:       os.Getenv("FB_ADMINSDK_CLIENT_X509_CERT_URL"),
	}

	credentialJSON, err := json.Marshal(credential)
	if err != nil {
		log.Fatalln(err)
	}

	sa := option.WithCredentialsJSON(credentialJSON)
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	// Authentication
	fa, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	authClient = fa

	// Cloud Messaging
	fm, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	messagingClient = fm

	log.Println("Firebase AdminSDK successfully initialized")
}

// Auth :: Get Authentication client
func Auth() *auth.Client {
	return authClient
}

// Messaging :: Get Cloud Messaging client
func Messaging() *messaging.Client {
	return messagingClient
}
