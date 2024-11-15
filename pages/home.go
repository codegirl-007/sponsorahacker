package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"sponsorahacker/db"
	"sponsorahacker/utils"
)

func Home(c *gin.Context) {

	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	isLoggedIn := utils.CheckIfLoggedIn(c)
	fmt.Println("isLoggedIn:", isLoggedIn)

	_, err := db.NewDbClient()
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":      "Sponsor a Hacker",
		"isLoggedIn": isLoggedIn,
	})
}
