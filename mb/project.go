package mb

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Project struct {
	ProjectID string `json:"uuid"`
	UserID string `json:"userID"`
	DisplayName string `json:"displayName"`
	CreatedAt time.Time `json:"createdAt"`
	Title string `json:"title"`
	Short string `json:"short"`
	LookingFor []string `form:"LookingFor[]"`
	Applicants []Applicant `json:"applicants"`
	Members []Member `json:"members"`
	Started bool
	StartedAt time.Time
}

type UserProject struct {
	Project *Project
	NeedsBackend bool
	NeedsFrontend bool
	NeedsInfra bool
}

type Applicant struct {
	ID string `json:"id"`
	DisplayName string `json:"displayName"`
	Which string `json:"which"` // frontend, backend or infra
}

type Member struct {
	ID string `json:"id"`
	DisplayName string `json:"displayName"`
}

type ProjectRequest struct {
	Title string `json:"title"`
	Short string `json:"short"`
	LookingFor []string `form:"LookingFor[]"`
}

func NewProject(userID string, pr *ProjectRequest) *Project {
	u := GetUserByID(userID)
	
	if u == nil {
		return nil
	}
	
	return &Project{
		ProjectID: uuid.New().String(),
		UserID: userID,
		DisplayName: u.DiscordUser.Username,
		CreatedAt: time.Now(),
		Title: pr.Title,
		Short: pr.Short,
		LookingFor: pr.LookingFor,
		Applicants: []Applicant{},
		Members: []Member{
			{
				ID: u.ID,
				DisplayName: u.DiscordUser.Username,
			},
		},
	}
}

func PostProject(c *gin.Context) {
	var projRequest ProjectRequest
	err := c.Bind(&projRequest)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, DataResponse(c, gin.H{
			"msg": "invalid request",
		}))
		return
	}
	
	userID, err := c.Cookie("user_id")
	if err != nil {
		fmt.Println("couldn't find user")
		return
	}
	
	newProj := NewProject(userID, &projRequest)
	
	// add to collection
	_, _, err = StoreClient.Collection("posts").Add(context.Background(), newProj)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	c.Redirect(http.StatusFound, "/users/" + userID + "/projects")
}

func UpdateProject(c *gin.Context) {
	var projRequest ProjectRequest
	err := c.Bind(&projRequest)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, DataResponse(c, gin.H{
			"msg": "invalid request",
		}))
		return
	}

	userID, err := c.Cookie("user_id")
	if err != nil {
		fmt.Println("couldn't find user")
		return
	}

	id := c.Param("id")
	
	posts, err := StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}

	for _, p := range posts {
		po := StoreClient.Collection("posts").Doc(p.Ref.ID)
		po.Update(context.Background(), []firestore.Update{
			{Path: "Title", Value: projRequest.Title},
			{Path: "Short", Value: projRequest.Short},
			{Path: "LookingFor", Value: projRequest.LookingFor},
		})
	}

	c.Redirect(http.StatusFound, "/users/" + userID + "/projects")
}

func DeleteProject(c *gin.Context) {

	userID, err := c.Cookie("user_id")
	if err != nil {
		fmt.Println("couldn't find user")
		return
	}

	id := c.Param("id")
	
	posts, err := StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}

	for _, p := range posts {
		var newP *Project
		p.DataTo(&newP)
		
		StoreClient.Collection("posts").Doc(p.Ref.ID).Delete(context.Background())
		break
	}

	c.Redirect(http.StatusFound, "/users/" + userID + "/projects")
}

func GetApprove(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	u := GetUserByID(userID)

	id := c.Param("id")
	applicantID := c.Param("applicantID")
	applicantUsername := c.Param("applicantUsername")

	posts, err := StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}

	for _, p := range posts {
		var newP *Project
		p.DataTo(&newP)

		newP.Members = append(newP.Members, Member{
			ID: applicantID,
			DisplayName: applicantUsername,
		})

		for i, v := range newP.Applicants {
			if v.ID == applicantID {
				newP.Applicants = append(newP.Applicants[:i], newP.Applicants[i+1:]...)
				break
			}
		}
		
		if newP.UserID == u.ID {
			po := StoreClient.Collection("posts").Doc(p.Ref.ID)
			_, err := po.Update(context.Background(), []firestore.Update{
				{Path: "Members", Value: newP.Members},
			})
			_, err = po.Update(context.Background(), []firestore.Update{
				{Path: "Applicants", Value: newP.Applicants},
			})
			if err != nil {
				// TODO
			}
		}
	}

	c.Redirect(http.StatusFound, "/users/" + userID + "/projects")
}

func GetDeny(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	u := GetUserByID(userID)

	id := c.Param("id")
	applicantID := c.Param("applicantID")

	posts, err := StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}

	for _, p := range posts {
		var newP *Project
		p.DataTo(&newP)

		for i, v := range newP.Applicants {
			if v.ID == applicantID {
				newP.Applicants = append(newP.Applicants[:i], newP.Applicants[i+1:]...)
				break
			}
		}

		if newP.UserID == u.ID {
			po := StoreClient.Collection("posts").Doc(p.Ref.ID)
			_, err = po.Update(context.Background(), []firestore.Update{
				{Path: "Applicants", Value: newP.Applicants},
			})
			if err != nil {
				// TODO
			}
		}
	}

	c.Redirect(http.StatusFound, "/users/" + userID + "/projects")
}

func GetRemove(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	u := GetUserByID(userID)

	id := c.Param("id")
	memberID := c.Param("memberID")

	posts, err := StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}

	for _, p := range posts {
		var newP *Project
		p.DataTo(&newP)

		for i, v := range newP.Members {
			if v.ID == memberID {
				newP.Members = append(newP.Members[:i], newP.Members[i+1:]...)
				break
			}
		}

		if newP.UserID == u.ID {
			po := StoreClient.Collection("posts").Doc(p.Ref.ID)
			_, err := po.Update(context.Background(), []firestore.Update{
				{Path: "Members", Value: newP.Members},
			})
			if err != nil {
				// TODO
			}
		}
	}

	c.Redirect(http.StatusFound, "/users/" + userID + "/projects")
}

func GetStart(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	
	id := c.Param("id")

	posts, err := StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
	if err != nil {
		return
	}

	for _, p := range posts {
		var newP *Project
		p.DataTo(&newP)

		if newP.UserID == userID && !newP.Started {
			newP.Started = true
			newP.StartedAt = time.Now()

			po := StoreClient.Collection("posts").Doc(p.Ref.ID)
			_, err := po.Update(context.Background(), []firestore.Update{
				{Path: "Started", Value: newP.Started},
				{Path: "StartedAt", Value: newP.StartedAt},
			})
			if err != nil {
				// TODO
			}
			
			break
		}
	}

	c.Redirect(http.StatusFound, "/users/" + userID + "/projects")
}
