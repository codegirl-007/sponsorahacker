package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
	db "sponsorahacker/db"
)

func Login(c *gin.Context) {
	providerName := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", providerName)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func Callback(c *gin.Context) {
	dbclient, err := db.NewDbClient()

	if err != nil {
		log.Println(err)
		return
	}

	sessionStore := NewSessionManager(dbclient)
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	if err != nil {
		log.Println("Error during user authentication:", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	err = sessionStore.CreateSession(user, c)

	if err != nil {
		log.Println("failed to create session in db:", err)
	}

	// For now, redirect to profile page after successful login
	c.Redirect(http.StatusTemporaryRedirect, "/welcome")
}

func Logout(c *gin.Context) {
	dbClient, err := db.NewDbClient()

	if err != nil {
		log.Fatal(err)
		return
	}

	sessionStore := NewSessionManager(dbClient)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = sessionStore.DeleteSession(c)

	if err != nil {
		log.Fatal("failed to delete session:", err)
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
