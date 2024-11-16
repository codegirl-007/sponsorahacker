package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sponsorahacker/utils"
)

func Welcome(c *gin.Context) {
	isLoggedIn := utils.CheckIfLoggedIn(c)
	fmt.Println("isLoggedIn:", isLoggedIn)
	c.HTML(http.StatusOK, "welcome.html", gin.H{
		"title":      "Sponsor A Hacker",
		"isLoggedIn": isLoggedIn,
	})
}
