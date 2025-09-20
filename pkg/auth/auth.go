package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"perfect-day/pkg/models"
	"perfect-day/pkg/storage"
	"time"
)

type AuthService struct {
	userStorage *storage.UserStorage
	sessions    map[string]*Session
}

type Session struct {
	ID       string
	Username string
	ExpiresAt time.Time
}

func NewAuthService(userStorage *storage.UserStorage) *AuthService {
	return &AuthService{
		userStorage: userStorage,
		sessions:    make(map[string]*Session),
	}
}

func (as *AuthService) Login(username string) (*models.User, *Session, error) {
	user, err := as.userStorage.Load(username)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %s", username)
	}

	session := &Session{
		ID:        generateSessionID(),
		Username:  username,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour sessions
	}

	as.sessions[session.ID] = session
	return user, session, nil
}

func (as *AuthService) ValidateSession(sessionID string) (*models.User, error) {
	session, exists := as.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("invalid session")
	}

	if time.Now().After(session.ExpiresAt) {
		delete(as.sessions, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	user, err := as.userStorage.Load(session.Username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", session.Username)
	}

	return user, nil
}

func (as *AuthService) Logout(sessionID string) {
	delete(as.sessions, sessionID)
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}