package mb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

type User struct {
	DiscordUser DiscordUser
	ID string
	Description string
	CreatedAt time.Time
	EnglishCreatedAt string
}

func NewUser(id, desc, discordJSON string) *User {
	return &User{
		DiscordUser: DiscordUserFromJSON(discordJSON),
		ID:          id,
		Description: desc,
		CreatedAt: time.Now(),
	}
}

func NewUserWithDiscordUser(id, desc string, du DiscordUser) *User {
	return &User{
		DiscordUser: du,
		ID:          id,
		Description: desc,
		CreatedAt: time.Now(),
	}
}

func GetUserByID(id string) *User {
	c := StoreClient.Collection("users")
	
	// check if user already exists
	all, err := c.Where("ID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return nil
	}
	
	if len(all) == 0 {
		return nil
	}
	
	for _, u := range all {
		var returnUser *User
		u.DataTo(&returnUser)
		return returnUser
	}
	
	return nil
}

func CreateOrUpdateUser(d DiscordUser) *User {
	u := GetUserByID(d.ID)
	if u != nil {
		return u
	}
	
	foundUser := NewUserWithDiscordUser(d.ID, "", d)
	foundUser.Save()
	return foundUser
}

// Save saves the user info into firestore
func (u *User) Save() {
	c := StoreClient.Collection("users")
	
	// check if user already exists
	all, err := c.Where("ID", "==", u.ID).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}
	
	// create if it doesn't exist
	if len(all) == 0 {
		_, _, err := c.Add(context.Background(), u)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	
	// update if it does
	for _, d := range all {
		uu := StoreClient.Collection("users").Doc(d.Ref.ID)
		uu.Update(context.Background(), []firestore.Update{
			{Path: "Description", Value: u.Description},
			{Path: "DiscordUser", Value: u.DiscordUser},
		})
	}
}

type DiscordUser struct {
	ID                   string `json:"id"`
	Username             string `json:"username"`
	Avatar               string `json:"avatar"`
	Discriminator        string `json:"discriminator"`
	PublicFlags          int    `json:"public_flags"`
	PremiumType          int    `json:"premium_type"`
	Flags                int    `json:"flags"`
	Banner               any    `json:"banner"`
	AccentColor          any    `json:"accent_color"`
	GlobalName           string `json:"global_name"`
	AvatarDecorationData any    `json:"avatar_decoration_data"`
	BannerColor          any    `json:"banner_color"`
	MfaEnabled           bool   `json:"mfa_enabled"`
	Locale               string `json:"locale"`
}

func DiscordUserFromJSON(j string) DiscordUser {
	var du DiscordUser
	err := json.Unmarshal([]byte(j), &du)
	if err != nil {
		return DiscordUser{}
	}

	return du
}

func DiscordUserToString(du DiscordUser) string {
	data, err := json.Marshal(&du)
	if err != nil {
		return ""
	}
	return string(data)
}
