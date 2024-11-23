package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"net/http"
	"sponsorahacker/auth"
)

func Welcome(c *gin.Context) {
	user, err := gothic.GetFromSession("user", c.Request)
	fmt.Println(user)

	if err != nil {
		fmt.Println("error finding session: ", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	userModel := auth.HydrateUser(user)
	fmt.Println(userModel.NickName)

	c.HTML(http.StatusOK, "welcome.html", gin.H{
		"title":      "Sponsor A Hacker",
		"isLoggedIn": true,
		"user":       userModel.NickName,
	})
}
