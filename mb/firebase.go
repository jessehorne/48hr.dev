package mb

import (
	"os"
)

type FirebaseClientCreds struct {
	APIKey string
	AuthDomain string
	ProjectID string
	StorageBucket string
	MessagingSenderID string
	AppID string
	MeasurementID string
}

func GetFirebaseClientCredentials() FirebaseClientCreds {
	return FirebaseClientCreds{
		APIKey: os.Getenv("FIREBASE_API_KEY"),
		AuthDomain: os.Getenv("FIREBASE_AUTH_DOMAIN"),
		ProjectID: os.Getenv("FIREBASE_PROJECT_ID"),
		StorageBucket: os.Getenv("FIREBASE_STORAGE_BUCKET"),
		MessagingSenderID: os.Getenv("FIREBASE_MESSAGING_SENDER_ID"),
		AppID: os.Getenv("FIREBASE_APP_ID"),
		MeasurementID: os.Getenv("FIREBASE_MEASUREMENT_ID"),
	}
}
