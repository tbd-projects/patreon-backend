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
)

type SessionsManager struct {
	sessionRep sessions.IRepository
}

func CreateSessionsManager(sessionRep sessions.IRepository) SessionsManager {
	return SessionsManager{sessionRep}
}

func (manager *SessionsManager) CheckSession(uniqID string) (models.Result, error) {
	userID, err := manager.sessionRep.GetUserId(uniqID)

	if err != nil {
		return models.Result{UserID: -1, UniqID: uniqID}, err
	}

	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return models.Result{UserID: -1, UniqID: uniqID}, err
	}
	return models.Result{UserID: intUserID, UniqID: uniqID}, nil
}

func generateUniqID(userID string) string {
	hash := md5.Sum(append([]byte(userID), uuid.NewV4().Bytes()...))
	return hex.EncodeToString(hash[:])
}

func (manager *SessionsManager) CreateSession(userID int64) (models.Result, error) {
	session := &models.Session{UserID: generateUniqID(fmt.Sprintf("%d", userID)),
		Expiration: durationStayCookies}

	err := manager.sessionRep.Set(session)

	if err != nil {
		return models.Result{UserID: -1}, err
	}

	return models.Result{UserID: userID, UniqID: session.UniqID}, nil
}

func (manager *SessionsManager) DeleteSession(uniqID string) error {
	session := &models.Session{UniqID: uniqID}
	err := manager.sessionRep.Del(session)
	return err
}
