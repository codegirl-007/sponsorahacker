package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
)

func Login(c *gin.Context) {
	providerName := c.Param("provider")

	// Begin the authentication process
	provider, err := goth.GetProvider(providerName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting provider: %s", err)
		return
	}

	state := uuid.New().String()

	session, err := provider.BeginAuth(state)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating auth url: %s", err)
		return
	}

	url, err := session.GetAuthURL()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting auth url: %s", err)
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func Callback(c *gin.Context) {
	sessionStore, err := NewSessionManager("libsql://sponsorahackersession-stephanie-gredell.turso.io")

	if err != nil {
		panic(err)
	}

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Println("Error during user authentication:", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	c.SetCookie("user_id", user.UserID, 3600, "/", "localhost", false, true)

	// For now, redirect to profile page after successful login
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Logout(c *gin.Context) {
	c.SetCookie("user_id", "", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
