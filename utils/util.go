package utils

import "github.com/gin-gonic/gin"

func CheckIfLoggedIn(c *gin.Context) bool {
	userID, err := c.Cookie("user_id")
	if err != nil || userID == "" {
		return false
	}

	return true
}
