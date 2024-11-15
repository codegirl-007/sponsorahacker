package pages

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Wishlist(c *gin.Context) {
	c.HTML(http.StatusOK, "wishlist.html", gin.H{
		"title": "Sponsor A Hacker",
	})
}
