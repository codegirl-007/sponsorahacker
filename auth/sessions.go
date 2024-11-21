package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/markbates/goth"
	"log"
	"net/url"
	"sponsorahacker/config"
	"sponsorahacker/db"
	"strconv"
	"time"
)

type SessionStore interface {
	CreateSession(*gin.Context) error
	GetSession(string) (Session, error)
	DeleteSession(string) error
}

type SessionManager struct {
	DB *db.DBClient
}

type Session struct {
	sessionId   string
	sessionData string
	createdOn   string
	modifiedOn  string
	expiresOn   string
}

type User struct {
	NickName        string `json:"NickName"`
	Provider        string `json:"Provider"`
	Provider_userid string `json:"UserID"`
}

var secureCookie *securecookie.SecureCookie

func NewSessionManager(db *db.DBClient) SessionManager {

	return SessionManager{DB: db}
}

func (s *SessionManager) CreateSession(user goth.User, c *gin.Context) error {
	sessionData, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("error marshalling user: %v", err)
	}

	auth := Session{
		sessionData: string(sessionData),
		createdOn:   time.Now().Format("20060102150405"),
		modifiedOn:  time.Now().Format("20060102150405"),
		expiresOn:   user.ExpiresAt.Format("20060102150405"),
	}

	result, err := s.DB.Exec(`
    INSERT INTO sessions (sessionData, createdOn, modifiedOn, expiresOn)
    VALUES (?, ?, ?, ?);`, auth.sessionData, auth.createdOn, auth.modifiedOn, auth.expiresOn)

	if result == nil {
		fmt.Printf("error getting result from database while creating the session: %v", err)
		return err
	}
	spuser, err := GetUserFromSession(auth)

	err = s.SaveUserIfNotExist(spuser)

	if err != nil {
		log.Printf("error creating session: %v", err)
		return err
	}

	sessionId, err := result.LastInsertId()

	if err != nil {
		fmt.Printf("error getting session id from database while creating the session: %v", err)
	}

	hash := []byte(config.GetEnvVar("COOKIE_HASH"))
	block := []byte(config.GetEnvVar("COOKIE_BLOCK"))
	secureCookie = securecookie.New(hash, block)
	cookieValue := map[string]string{
		"sessionId": strconv.Itoa(int(sessionId)),
	}

	encoded, err := secureCookie.Encode("_session", cookieValue)

	if err != nil {
		fmt.Printf("error encoding cookie value: %v", err)
		return err
	}

	c.SetCookie("_session", encoded, 3600, "/", "localhost", false, true)

	return nil
}

func (s *SessionManager) GetSession(session sessions.Session) (Session, error) {
	// query for one row
	result, err := s.DB.Query(`SELECT sessionData FROM sessions WHERE sessionId=$1 LIMIT 1`, session.ID())

	// if err, then return an empty struct
	if err != nil {
		return Session{}, err
	}

	// else go through the results and create a Session
	for result.Next() {
		var s Session

		// unless there is an error, of course, then return an empty struct
		if err := result.StructScan(&s); err != nil {
			return Session{}, err
		}

		return s, nil
	}

	// if we get nothing, well, we go nothing
	return Session{}, nil
}

func GetUserFromSession(session Session) (User, error) {
	var user User

	fmt.Println(session.sessionData)
	err := json.Unmarshal([]byte(session.sessionData), &user)

	if err != nil {
		fmt.Println("error unmarshalling session data for user", err)

		return User{}, err
	}

	return user, nil
}

func (s *SessionManager) SaveUserIfNotExist(user User) error {
	rows, err := s.DB.Query(`SELECT * FROM users WHERE provider_userid = ? AND provider = ? LIMIT 1`, user.Provider_userid, user.Provider)

	if err != nil {
		fmt.Printf("error getting user from database: %v", err)
		return err
	}

	exists := rows.Next()

	if !exists {
		_, err = s.DB.Exec(`
	  INSERT INTO users (provider_userid, provider, username)
	  VALUES (?, ?, ?);`, user.Provider_userid, user.Provider, user.NickName)

		if err != nil {
			fmt.Printf("error inserting user: %v", err)
			return err
		}
	}

	return nil
}

func (s *SessionManager) DeleteSession(c *gin.Context) error {
	if cookie, err := c.Request.Cookie("_session"); err == nil {
		value := make(map[string]string)

		cookieValue, _ := url.QueryUnescape(cookie.Value)

		hash := []byte(config.GetEnvVar("COOKIE_HASH"))
		block := []byte(config.GetEnvVar("COOKIE_BLOCK"))
		secureCookie := securecookie.New(hash, block)

		err = secureCookie.Decode("_session", cookieValue, &value)

		if err != nil {
			fmt.Printf("error decoding cookie value: %v", err)
			return err
		}

		sessionId := value["sessionId"]
		sessionIdInt, err := strconv.Atoi(sessionId)

		if err != nil {
			fmt.Printf("error converting sessionId to int: %v", err)
			return err
		}

		_, err = s.DB.Exec(`DELETE FROM sessions WHERE ID = ?`, sessionIdInt)

		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
	}

	c.SetCookie("_session", "", -1, "/", "localhost", false, true)

	return nil
}
