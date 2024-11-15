package pages

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", gin.H{
		"title": "Sponsor A Hacker",
	})
}
