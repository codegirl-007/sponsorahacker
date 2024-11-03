package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"sponsorahacker/config"
	"sponsorahacker/db"
)

type SessionManager interface {
	SetSession(username string, c *gin.Context) error
	GetSession(c *gin.Context) (string, error)
}

type SessionStore struct {
	SessionDB db.Database
	Store     *sessions.CookieStore
}

func NewSessionManager(dbUrl string) (*SessionStore, error) {
	db, err := db.NewDbClient(dbUrl)
	if err != nil {
		return nil, err
	}

	// Create sessions table if not exist
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL UNIQUE,
		data BLOB NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`)

	if err != nil {
		return nil, err
	}

	sessionSecret := config.GetEnvVar("SESSION_SECRET")

	store := sessions.NewCookieStore([]byte(sessionSecret))
	return &SessionStore{db, store}, nil
}

func (s *SessionStore) SetSession(username string, c *gin.Context) error {
	session, err := s.Store.Get(c.Request, "session")
	if err != nil {
		return err
	}

	session.Values["username"] = username
	if err := session.Save(c.Request, c.Writer); err != nil {
		return err
	}

	return s.saveSessionToDB(session)
}

func (s *SessionStore) GetSession(c *gin.Context) (string, error) {
	session, err := s.Store.Get(c.Request, "session")
	if err != nil {
		return "", err
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		return "", fmt.Errorf("username not found in session")
	}

	return username, nil
}

func (s *SessionStore) DeleteSession(c *gin.Context) error {
	session, err := s.Store.Get(c.Request, "session")
	if err != nil {
		return err
	}

	session.Values["username"] = make(map[interface{}]interface{})

	if err := session.Save(c.Request, c.Writer); err != nil {
		return err
	}

	_, err = s.SessionDB.Exec("DELETE FROM sessions WHERE session_id = ?", session.ID)

	return err
}

func (s *SessionStore) saveSessionToDB(session *sessions.Session) error {
	data, err := json.Marshal(session.Values)
	if err != nil {
		return err
	}

	_, err = s.SessionDB.Exec(`
	INSERT INTO sessions (session_id, data, updated_at)
	VALUES (?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(session_id) DO UPDATE SET data = ?, updated_at = CURRENT_TIMESTAMP
	`, session.ID, data, data)

	return err
}
