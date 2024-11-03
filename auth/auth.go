package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
	"sponsorahacker/config"
)

func Login(c *gin.Context) {
	providerName := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", providerName)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func Callback(c *gin.Context) {
	sessionStore, err := NewSessionManager(config.GetEnvVar("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	if err != nil {
		log.Println("Error during user authentication:", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	c.SetCookie("user_id", user.UserID, 3600, "/", "localhost", false, true)
	err = sessionStore.SetSession(user.Name, c)
	if err != nil {
		log.Println("failed to set session:", err)
	}

	// For now, redirect to profile page after successful login
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Logout(c *gin.Context) {
	sessionStore, err := NewSessionManager(config.GetEnvVar("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	c.SetCookie("user_id", "", -1, "/", "localhost", false, true)
	err = sessionStore.DeleteSession(c)
	if err != nil {
		log.Println("failed to delete session:", err)
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
