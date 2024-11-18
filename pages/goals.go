package pages

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sponsorahacker/db"
	"sponsorahacker/utils"
	"strconv"
	"strings"
	"time"
)

type Goal struct {
	Name          string  `form:"item-name" db:"name" binding:"required"`
	FundingAmount float64 `form:"funding-amount" db:"funding_amount" binding:"required,numeric"`
	Description   string  `form:"item-description" db:"description" binding:"required"`
	LearnMoreURL  *string `form:"item-url" db:"learn_more_url"` // Optional field
	CreatedAt     string  `db:"created_at"`
	UpdatedAt     string  `db:"updated_at"`
}

func Goals(c *gin.Context) {
	isLoggedIn := utils.CheckIfLoggedIn(c)

	c.HTML(http.StatusOK, "goals.html", gin.H{
		"title":      "Sponsor A Hacker",
		"isLoggedIn": isLoggedIn,
	})
}

func CreateGoal(c *gin.Context) {
	dbClient, err := db.NewDbClient()

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		return
	}

	isLoggedIn := utils.CheckIfLoggedIn(c)

	name := strings.TrimSpace(c.PostForm("item-name"))
	fundingAmountStr := strings.TrimSpace(c.PostForm("funding-amount"))
	description := strings.TrimSpace(c.PostForm("item-description"))
	learnMoreURL := strings.TrimSpace(c.PostForm("item-url"))

	if name == "" && fundingAmountStr == "" && description == "" && learnMoreURL == "" {
		fieldError := fmt.Errorf("missing a required field. Please check and make sure all fields are filled out")
		c.HTML(http.StatusOK, "goals.html", gin.H{
			"title":      "Sponsor A Hacker",
			"isLoggedIn": isLoggedIn,
			"error":      fieldError,
		})
	}

	createdDate := time.Now()
	updatedDate := time.Now()
	fundingAmount, err := strconv.ParseFloat(fundingAmountStr, 64)

	if err != nil {
		fieldError := fmt.Errorf("invalid funding amount")
		c.HTML(http.StatusOK, "goals.html", gin.H{
			"title":      "Sponsor A Hacker",
			"isLoggedIn": isLoggedIn,
			"error":      fieldError,
		})
	}

	_, err = dbClient.Exec(`INSERT INTO goals (name, description, learn_more_url, funding_amount, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);`, name, description, learnMoreURL, fundingAmount, createdDate, updatedDate)

	if err != nil {
		fmt.Printf("Error inserting goal: %v", err)
		c.HTML(http.StatusOK, "goals.html", gin.H{
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
