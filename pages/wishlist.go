package pages

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sponsorahacker/utils"
)

func Goals(c *gin.Context) {
	isLoggedIn := utils.CheckIfLoggedIn(c)

	c.HTML(http.StatusOK, "goals.html", gin.H{
		"title":      "Sponsor A Hacker",
		"isLoggedIn": isLoggedIn,
	})
}
