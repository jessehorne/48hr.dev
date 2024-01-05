package mb

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
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
	NeedsBackend bool `json:"needsBackend"`
	NeedsFrontend bool `json:"needsFrontend"`
	NeedsInfra bool `json:"needsInfra"`
	Applicants []Applicant `json:"applicants"`
	Members []string `json:"members"`
}

type Applicant struct {
	ID string `json:"id"`
	DisplayName string `json:"displayName"`
	Which string `json:"which"` // frontend, backend or infra
}

type ProjectRequest struct {
	Title string `json:"title"`
	Short string `json:"short"`
	NeedsBackend bool `json:"needsBackend"`
	NeedsFrontend bool `json:"needsFrontend"`
	NeedsInfra bool `json:"needsInfra"`
}

func NewProject(userID string, pr *ProjectRequest) *Project {
	u, err := AuthClient.GetUser(context.Background(), userID)
	if err != nil {
		return nil
	}
	
	
	return &Project{
		ProjectID: uuid.New().String(),
		UserID: userID,
		DisplayName: u.DisplayName,
		CreatedAt: time.Now(),
		Title: pr.Title,
		Short: pr.Short,
		NeedsBackend: pr.NeedsBackend,
		NeedsFrontend: pr.NeedsFrontend,
		NeedsInfra: pr.NeedsInfra,
		Applicants: []Applicant{},
		Members: []string{userID},
	}
}

func PostProject(c *gin.Context) {
	var projRequest ProjectRequest
	err := c.Bind(&projRequest)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid request",
		})
		return
	}
	
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
	
	newProj := NewProject(userID, &projRequest)
	
	// add to collection
	StoreClient.Collection("posts").Add(context.Background(), newProj)
	
	c.JSON(http.StatusOK, gin.H{})
}
