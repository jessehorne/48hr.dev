package mb

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func GetApply(c *gin.Context) {
	// get user id from request
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "not authenticated",
		})
		return
	}

	t := token.(*auth.Token)
	if t == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid token",
		})
		return
	}

	userID := t.UID

	u, err := AuthClient.GetUser(context.Background(), userID)
	if err != nil {
		return
	}

	which := c.Param("which")
	id := c.Param("id")

	newApplicant := Applicant{
		ID:          userID,
		DisplayName: u.DisplayName,
		Which: which,
	}

	posts, err := StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}

	for _, p := range posts {
		var newP *Project
		p.DataTo(&newP)
		newP.Applicants = append(newP.Applicants, newApplicant)

		po := StoreClient.Collection("posts").Doc(p.Ref.ID)
		_, err := po.Update(context.Background(), []firestore.Update{
			{Path: "Applicants", Value: newP.Applicants},
		})
		if err != nil {
			// TODO
		}
	}

	c.Redirect(http.StatusFound, "/")
}
