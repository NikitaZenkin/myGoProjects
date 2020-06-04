package session

import (
	"database/sql"
	"net/http"
	"rclone/pkg/system"
	"time"
)

type SessionsManager struct {
	DB *sql.DB
}

func NewSesManager(db *sql.DB) SessionsManager {
	return SessionsManager{DB: db}
}

func (sm *SessionsManager) CreateSession(w http.ResponseWriter, userID string, userName string) error {
	sessionId := system.NewId()

	sm.DB.Exec("DELETE FROM sessions WHERE `userID` = ?",
		userID,
	)
	_, err := sm.DB.Exec(
		"INSERT INTO sessions (`ID`,`userID`, `userName`) VALUES (?, ?, ?)",
		sessionId,
		userID,
		userName,
	)
	if err != nil {
		return err
	}
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sessionId,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	return nil
}

func (sm *SessionsManager) CheckSession(r *http.Request) (*Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, err
	}
	sessionId := sessionCookie.Value
	session := &Session{}
	err = sm.DB.
		QueryRow("SELECT ID, userID, userName FROM sessions WHERE ID = ?", sessionId).
		Scan(&session.ID, &session.UserId, &session.UserName)
	if err != nil {
		return nil, err
	}
	return session, nil

}
