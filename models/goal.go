package models

import (
	"fmt"
	"log"
	"time"

	"sponsorahacker/db"
)

type Goal struct {
	id            int
	Name          string    `form:"item-name" db:"name" binding:"required"`
	FundingAmount float64   `form:"funding-amount" db:"funding_amount" binding:"required,numeric"`
	Description   string    `form:"item-description" db:"description" binding:"required"`
	LearnMoreURL  string    `form:"item-url" db:"learn_more_url"` // Optional field
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	CreatedBy     int       `db:"created_by"`
}

func (g *Goal) CreateGoal() error {
	dbClient, err := db.NewDbClient()

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		return err
	}
	_, err = dbClient.Exec(`INSERT INTO goals (name, description, learn_more_url, funding_amount, created_by, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);`, g.Name, g.Description, g.LearnMoreURL, g.FundingAmount, g.CreatedBy, g.CreatedAt, g.UpdatedAt)

	if err != nil {
		log.Fatalf("Could not create goal: %v", err)
	}
	return nil
}

func GetGoals(user int) ([]Goal, error) {
	var goals []Goal

	dbClient, err := db.NewDbClient()

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	result, err := dbClient.Query(`SELECT name, description, learn_more_url, funding_amount, created_by, created_at, updated_at FROM goals where created_by = ?`, user)

	if err != nil {
		log.Fatalf("Could not get goals: %v", err)
	}
	for result.Next() {
		var goal Goal
		err = result.StructScan(&goal)

		fmt.Println(goal.Name)

		if err != nil {
			log.Fatalf("Could not read goal: %v", err)
			return nil, err
		}
		goals = append(goals, goal)
	}

	return goals, nil
}

func GetGoal(id int) (*Goal, error) {
	dbClient, err := db.NewDbClient()
	result, err := dbClient.Query(`SELECT name, description, learn_more_url, funding_amount, created_by, created_at, updated_at FROM goals where id = ?`, id)

	if err != nil {
		log.Fatalf("Could not read goal: %v", err)
		return nil, err
	}

	var goal Goal
	err = result.StructScan(&goal)

	return &goal, nil
}
