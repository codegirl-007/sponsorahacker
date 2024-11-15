package utils

import "github.com/gin-gonic/gin"

func CheckIfLoggedIn(c *gin.Context) bool {
	cookie, err := c.Cookie("_session")
	if err != nil || cookie == "" {
		return false
	}

	return true
}
