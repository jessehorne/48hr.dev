package mb

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/go-playground/validator/v10"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

var Context context.Context
var App *firebase.App
var AuthClient *auth.Client
var StoreClient *firestore.Client
var Validator *validator.Validate
var DebugMode bool
var State string
var OAuthConf *oauth2.Config

func InitApp(debug bool) {
	time.Local = time.UTC
	
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
	State, _ = generateRandomString(32)

	// Create an oauth2 config.
	// Ensure you add the redirect url in the application's oauth2 settings
	// in the discord devs page.
	conf := &oauth2.Config{
		RedirectURL: "https://localhost:8080/auth/callback",
		// This next 2 lines must be edited before running this.
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		Scopes:       []string{discord.ScopeIdentify},
		Endpoint:     discord.Endpoint,
	}
	OAuthConf = conf
}
