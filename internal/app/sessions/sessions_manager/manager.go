package sessions_manager

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/satori/go.uuid"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/models"
	"strconv"
)

const (
	oneDayInMillisecond = 86400
	durationStayCookies = oneDayInMillisecond * 2
	UnknownUser         = -1
)

type SessionsManager struct {
	sessionRep sessions.Repository
}

func CreateSessionsManager(sessionRep sessions.Repository) SessionsManager {
	return SessionsManager{sessionRep}
}

func (manager *SessionsManager) CheckSession(uniqID string) (models.Result, error) {
	userID, err := manager.sessionRep.GetUserId(uniqID)

	if err != nil {
		return models.Result{UserID: UnknownUser, UniqID: uniqID}, err
	}

	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return models.Result{UserID: UnknownUser, UniqID: uniqID}, err
	}
	return models.Result{UserID: intUserID, UniqID: uniqID}, nil
}

func generateUniqID(userID string) string {
	hash := md5.Sum(append([]byte(userID), uuid.NewV4().Bytes()...))
	return hex.EncodeToString(hash[:])
}

func (manager *SessionsManager) CreateSession(userID int64) (models.Result, error) {
	stringUserID := fmt.Sprintf("%d", userID)
	session := &models.Session{UniqID: generateUniqID(stringUserID), UserID: stringUserID,
		Expiration: durationStayCookies}

	err := manager.sessionRep.Set(session)

	if err != nil {
		return models.Result{UserID: UnknownUser}, err
	}

	return models.Result{UserID: userID, UniqID: session.UniqID}, nil
}

func (manager *SessionsManager) DeleteSession(uniqID string) error {
	session := &models.Session{UniqID: uniqID}
	err := manager.sessionRep.Del(session)
	return err
}
