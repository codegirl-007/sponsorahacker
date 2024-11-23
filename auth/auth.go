package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
	"sponsorahacker/db"
	"strconv"
)

func Login(c *gin.Context) {
	providerName := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", providerName)
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func Callback(c *gin.Context) {
	// complete the authentication process
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	// do the error checking
	if err != nil {
		log.Println("Error during user authentication:", err)
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	// serialize the data
	userData, err := json.Marshal(user)

	// check error for serialization
	if err != nil {
		log.Println("Error serializing user to json: ", err)
	}

	dbClient, err := db.NewDbClient()

	if err != nil {
		log.Println("Error connecting to db: ", err)
	}

	// new stuff
	selectQuery := `SELECT id FROM users WHERE provider_userid = ? AND provider = ?`

	rows, err := dbClient.Query(selectQuery, user.UserID, user.Provider)
	if err != nil {
		log.Fatalln("Error executing query:", err)
	}

	defer rows.Close()

	var userId int64

	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			log.Fatalln("Error scanning row:", err)
		}
		fmt.Println("User already exists with ID:", userId)
	} else {
		insertQuery := `INSERT INTO users (provider_userid, provider, nickname) VALUES (?, ?, ?)`
		result, err := dbClient.Exec(insertQuery, user.UserID, user.Provider, user.NickName)
		if err != nil {
			log.Fatalln("Error inserting user:", err)
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			log.Fatalln("Error getting last insert ID:", err)
		}

		userId = lastInsertID
	}
	// old stuff

	// unserialize so we can add more things
	var userMap map[string]interface{}
	if err := json.Unmarshal(userData, &userMap); err != nil {
		log.Println("Error unmarshalling user JSON:", err)
		return
	}
	userMap["Sid"] = strconv.FormatInt(userId, 10)
	updatedUserData, err := json.Marshal(userMap)

	if err != nil {
		log.Println("Error serializing user JSON:", err)
	}

	// store data in a session
	err = gothic.StoreInSession("user", string(updatedUserData), c.Request, c.Writer)

	// more error checking
	if err != nil {
		fmt.Println("Error saving user to session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// For now, redirect to welcome page
	c.Redirect(http.StatusTemporaryRedirect, "/welcome")
}

func Logout(c *gin.Context) {
	session, err := gothic.Store.Get(c.Request, "_gothic-session")
	if err != nil {
		fmt.Println("Error retrieving session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve session"})
		return
	}

	// Clear the session data
	session.Values = make(map[interface{}]interface{})

	// Save the empty session
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		fmt.Println("Error saving session: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
