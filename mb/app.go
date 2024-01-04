package mb

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/go-playground/validator/v10"
	"google.golang.org/api/option"
)

var Context context.Context
var App *firebase.App
var AuthClient *auth.Client
var StoreClient *firestore.Client
var Validator *validator.Validate
var DebugMode bool

func InitApp(debug bool) {
	DebugMode = debug
	
	Context = context.Background()

	// initialize firebase
	opt := option.WithCredentialsFile("./key.json")

	app, err := firebase.NewApp(Context, nil, opt)
	if err != nil {
		panic(err)
	}

	authClient, err := app.Auth(Context)
	if err != nil {
		panic(err)
	}

	storeClient, err := app.Firestore(Context)
	if err != nil {
		panic(err)
	}

	App = app
	AuthClient = authClient
	StoreClient = storeClient
	Validator = validator.New()
}
