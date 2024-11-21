package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"sponsorahacker/models"
	"sponsorahacker/utils"
	"strconv"
	"strings"
	"time"
)

func Goals(c *gin.Context) {
	isLoggedIn := utils.CheckIfLoggedIn(c)

	c.HTML(http.StatusOK, "goals.html", gin.H{
		"title":      "Sponsor A Hacker",
		"isLoggedIn": isLoggedIn,
	})
}

func CreateGoal(c *gin.Context) {
	isLoggedIn := utils.CheckIfLoggedIn(c)
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

	_, err := parseUrl(learnMoreURL)

	if isEmpty(name) && isEmpty(fundingAmountStr) && isEmpty(description) && isEmpty(learnMoreURL) && err != nil {
		fieldError := fmt.Errorf("missing a required field. Please check and make sure all fields are filled out")
		c.HTML(http.StatusNotAcceptable, "goals.html", gin.H{
			"title":      "Sponsor A Hacker",
			"isLoggedIn": isLoggedIn,
			"error":      fieldError,
		})
	}

	fundingAmount, err := parseFloat(strings.Replace(fundingAmountStr, ",", "", -1), 64)

	if err != nil {
		fieldError := fmt.Errorf("invalid funding amount")
		c.HTML(http.StatusInternalServerError, "goals.html", gin.H{
			"title":      "Sponsor A Hacker",
			"isLoggedIn": isLoggedIn,
			"error":      fieldError,
		})
	}

	createdDate := time.Now()
	updatedDate := time.Now()

	goal := models.Goal{
		Name:          name,
		FundingAmount: fundingAmount,
		Description:   description,
		LearnMoreURL:  learnMoreURL,
		CreatedAt:     createdDate,
		UpdatedAt:     updatedDate,
	}

	err = goal.CreateGoal()

	if err != nil {
		fmt.Printf("Error inserting goal: %v", err)
		c.HTML(http.StatusInternalServerError, "goals.html", gin.H{
			"title":      "Sponsor A Hacker",
			"isLoggedIn": isLoggedIn,
			"error":      err,
		})
	}

	c.HTML(http.StatusOK, "goals.html", gin.H{
		"title":      "Sponsor A Hacker",
		"isLoggedIn": isLoggedIn,
	})
}
