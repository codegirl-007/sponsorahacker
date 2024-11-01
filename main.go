package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"log"
	"net/http"
	"sponsorahacker/db"
)

// route to authenticate
// route for log
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		isLoggedIn := checkIfLoggedIn(c)
		fmt.Println("isLoggedIn:", isLoggedIn)

		_, err := db.NewDbClient("libsql://sponsorahacker-stephanie-gredell.turso.io")
		if err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Sponsor a Hacker",
			"isLoggedIn": isLoggedIn,
		})
	})

	r.GET("/apply", func(c *gin.Context) {
		c.HTML(http.StatusOK, "apply.html", gin.H{
			"title": "Apply",
		})
	})

	goth.UseProviders(
		twitch.New("3bozqr9birw9qcpdmoeo4my9vjuvoa", "bidyehx9w9csdqqtx95ndpa7zwrjir", "http://localhost:8080/auth/twitch/callback"),
		google.New("695096681313-07o8ocke1ifpuim1c9n4k1uacbdvrb4j.apps.googleusercontent.com", "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", "http://localhost:8080/auth/google/callback"),
		twitter.New("WSuqS3zaiCFEQOVwXXUq68Duj", "xuIh2dLwXO7P2UeGiBTDUmAfTPqkRdPx4PxwCP88Mj9oSXBg56", "http://localhost:8080/auth/twitter/callback"),
	)

	r.Use(func(c *gin.Context) {
		gothic.GetProviderName = func(req *http.Request) (string, error) {
			return "twitch", nil
		}
		c.Next()
	})

	// do the authentication thing
	r.GET("/auth/login/:provider", func(c *gin.Context) {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	})

	r.GET("/auth/:provider/callback", func(c *gin.Context) {
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			log.Println("Error during user authentication:", err)
			c.Redirect(http.StatusTemporaryRedirect, "/homepage")
			return
		}
		c.SetCookie("user_id", user.UserID, 3600, "/", "localhost", false, true)

		// For now, redirect to profile page after successful login
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	// do the logout thang
	r.GET("/auth/logout/:provider", func(c *gin.Context) {
		c.SetCookie("user_id", "", -1, "/", "localhost", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	err := r.Run(":8080")

	if err != nil {
		panic(err)
	}
}

func checkIfLoggedIn(c *gin.Context) bool {
	userID, err := c.Cookie("user_id")
	if err != nil || userID == "" {
		return false
	}

	return true
}
