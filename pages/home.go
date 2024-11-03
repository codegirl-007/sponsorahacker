package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sponsorahacker/db"
	"sponsorahacker/twitch"
	"sponsorahacker/utils"
)

func Home(c *gin.Context) {

	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	isLoggedIn := utils.CheckIfLoggedIn(c)
	fmt.Println("isLoggedIn:", isLoggedIn)

	_, err := db.NewDbClient("libsql://sponsorahacker-stephanie-gredell.turso.io")
	if err != nil {
		log.Fatal(err)
	}

	client, twitchErr := twitch.NewTwitchClient(os.Getenv("TWITCH_CLIENT_ID"), os.Getenv("TWITCH_CLIENT_SECRET"))

	if twitchErr != nil {
		log.Fatal(twitchErr)
	}

	client.GetUser()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "Sponsor a Hacker",
		"isLoggedIn": isLoggedIn,
	})
}
