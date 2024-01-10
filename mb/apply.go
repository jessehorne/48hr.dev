package mb

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func GetApply(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	u := GetUserByID(userID)

	which := c.Param("which")
	id := c.Param("id")

	newApplicant := Applicant{
		ID:          userID,
		DisplayName: u.DiscordUser.Username,
		Which:       which,
	}

	posts, err := StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}

	for _, p := range posts {
		var newP *Project
		p.DataTo(&newP)

		// only add to list of applicants if it isn't already in there
		var found bool
		for _, a := range newP.Applicants {
			if a.ID == userID && a.Which == which {
				found = true
				break
			}
		}

		if !found {
			newP.Applicants = append(newP.Applicants, newApplicant)

			po := StoreClient.Collection("posts").Doc(p.Ref.ID)
			_, err := po.Update(context.Background(), []firestore.Update{
				{Path: "Applicants", Value: newP.Applicants},
			})
			if err != nil {
				// TODO
			}
		}
	}

	c.Redirect(http.StatusFound, "/")
}

func GetDisable(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	u := GetUserByID(userID)

	which := c.Param("which")
	id := c.Param("id")

	posts, err := StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}

	for _, p := range posts {
		var newP *Project
		p.DataTo(&newP)

		for i, v := range newP.LookingFor {
			if v == which {
				newP.LookingFor = append(newP.LookingFor[:i], newP.LookingFor[i+1:]...)
				break
			}
		}

		if newP.UserID == u.ID {
			po := StoreClient.Collection("posts").Doc(p.Ref.ID)
			_, err := po.Update(context.Background(), []firestore.Update{
				{Path: "LookingFor", Value: newP.LookingFor},
			})
			if err != nil {
				// TODO
			}
		}
	}

	c.Redirect(http.StatusFound, "/")
}
