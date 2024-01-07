package mb

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xeonx/timeago"
)

func DataResponse(c *gin.Context, data gin.H) gin.H {
	var isAuthed bool
	authed, _ := c.Cookie("is_authed")
	if authed == "true" {
		isAuthed = true
	}

	data["authed"] = isAuthed

	id, err := c.Cookie("user_id")
	fmt.Println("USER ID: ", id)
	if err == nil {
		u := GetUserByID(id)
		if u != nil {
			data["UserID"] = u.ID
			data["UserDiscordID"] = u.DiscordUser.ID
			data["UserDisplayName"] = u.DiscordUser.Username
			data["UserDescription"] = u.Description
			data["UserCreatedAt"] = u.CreatedAt
			data["EnglishCreatedAt"] = timeago.English.Format(u.CreatedAt)
		} else {
			fmt.Println("couldn't get user by id")
		}
	} else {
		fmt.Println("couldn't find user")
	}

	return data
}
