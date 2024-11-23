package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"net/http"
	"sponsorahacker/auth"
)

func Home(c *gin.Context) {
	user, err := gothic.GetFromSession("user", c.Request)

	if err != nil {
		fmt.Println("error checking login", err)
	}

	auth.HydrateUser(user)

	if err != nil {
		fmt.Println("error checking login", err)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Sponsor a Hacker",
			"isLoggedIn": false,
		})

		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/welcome")
}
