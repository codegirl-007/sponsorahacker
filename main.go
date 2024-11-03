package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/twitch"
	"log"
	"os"
	"sponsorahacker/auth"
	"sponsorahacker/pages"
)

// route to authenticate
// route for log
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	goth.UseProviders(
		twitch.New(os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_SECRET"), "http://localhost:8080/auth/twitch/callback"),
		// google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_SECRET"), "http://localhost:8080/auth/google/callback"),
		// twitter.New(os.Getenv("TWITTER_CLIENT_ID"), os.Getenv("TWITTER_SECRET"), "http://localhost:8080/auth/twitter/callback"),
	)

	// auth routes
	r.GET("/auth/login/:provider", auth.Login)
	r.GET("/auth/:provider/callback", auth.Callback)
	r.GET("/auth/logout/:provider", auth.Logout)

	// pages routes
	r.GET("/", pages.Home)
	r.GET("/login", pages.Login)

	runErr := r.Run(":8080")

	if runErr != nil {
		panic(runErr)
	}
}
