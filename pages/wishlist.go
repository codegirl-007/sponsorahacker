package pages

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Goals(c *gin.Context) {
	c.HTML(http.StatusOK, "goals.html", gin.H{
		"title": "Sponsor A Hacker",
	})
}
