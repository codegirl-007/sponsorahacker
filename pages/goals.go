package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
	"net/url"
	"sponsorahacker/auth"
	"sponsorahacker/models"
	"strconv"
	"strings"
	"time"
)

func Goals(c *gin.Context) {
	user, err := gothic.GetFromSession("user", c.Request)

	if err != nil {
		fmt.Println("error finding session: ", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	userModel := auth.HydrateUser(user)

	goals, err := models.GetGoals(userModel.GetSid())

	if err != nil {
		fmt.Println("error getting goals: ", err)
	}

	fmt.Println("sid: ", userModel.GetSid())
	fmt.Println("goals:", goals)

	c.HTML(http.StatusOK, "goals.html", gin.H{
		"title":      "Sponsor A Hacker",
		"isLoggedIn": true,
		"user":       userModel.NickName,
		"goals":      goals,
	})
}

func Goal(c *gin.Context) {
	user, err := gothic.GetFromSession("user", c.Request)

	if err != nil {
		fmt.Println("error finding session: ", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	goalId, err := strconv.Atoi(c.Param("goalId"))

	if err != nil {
		log.Println("error parsing goal id: ", err)
	}

	userModel := auth.HydrateUser(user)
	goal, err := models.GetGoal(goalId)

	if err != nil {
		fmt.Println("error getting goals: ", err)
	}

	c.HTML(http.StatusOK, "goal.html", gin.H{
		"title":      "Sponsor A Hacker",
		"isLoggedIn": true,
		"user":       userModel.NickName,
		"goal":       goal,
	})
}

func CreateGoal(c *gin.Context) {
	user, err := gothic.GetFromSession("user", c.Request)

	if err != nil {
		fmt.Println("error finding session: ", err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	userModel := auth.HydrateUser(user)

	trim := strings.TrimSpace
	parseFloat := strconv.ParseFloat
	parseUrl := url.Parse

	isEmpty := func(s string) bool {
		return trim(s) == ""
	}

	name := trim(c.PostForm("item-name"))
	fundingAmountStr := trim(c.PostForm("funding-amount"))
	description := trim(c.PostForm("item-description"))
	learnMoreURL := trim(c.PostForm("item-url"))

	_, err = parseUrl(learnMoreURL)

	if isEmpty(name) && isEmpty(fundingAmountStr) && isEmpty(description) && isEmpty(learnMoreURL) && err != nil {
		fieldError := fmt.Errorf("missing a required field. Please check and make sure all fields are filled out")
		c.HTML(http.StatusNotAcceptable, "goals.html", gin.H{
			"title":      "Sponsor A Hacker",
			"isLoggedIn": true,
			"error":      fieldError,
		})
	}

	fundingAmount, err := parseFloat(strings.Replace(fundingAmountStr, ",", "", -1), 64)

	if err != nil {
		fieldError := fmt.Errorf("invalid funding amount")
		c.HTML(http.StatusInternalServerError, "goals.html", gin.H{
			"title":      "Sponsor A Hacker",
			"isLoggedIn": true,
			"error":      fieldError,
		})
	}

	createdDate := time.Now()
	updatedDate := time.Now()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "goals.html", gin.H{})
	}

	fmt.Println("sid is : ", userModel.GetSid())
	goal := models.Goal{
		Name:          name,
		FundingAmount: fundingAmount,
		Description:   description,
		LearnMoreURL:  learnMoreURL,
		CreatedAt:     createdDate,
		CreatedBy:     userModel.GetSid(),
		UpdatedAt:     updatedDate,
	}

	err = goal.CreateGoal()

	if err != nil {
		fmt.Printf("Error inserting goal: %v", err)
		c.HTML(http.StatusInternalServerError, "goals.html", gin.H{
			"title":      "Sponsor A Hacker",
			"isLoggedIn": true,
			"error":      err,
		})
	}

	c.HTML(http.StatusOK, "goals.html", gin.H{
		"title":      "Sponsor A Hacker",
		"isLoggedIn": true,
	})
}
