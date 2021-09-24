package sessions_manager

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/models"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

const (
	oneDayInMillisecond = 86400
	expiredCookiesTime  = oneDayInMillisecond * 2
	UnknownUser         = -1
)

type SessionManager struct {
	sessionRepository sessions.SessionRepository
}

func NewSessionManager(sessionRep sessions.SessionRepository) *SessionManager {
	return &SessionManager{
		sessionRepository: sessionRep,
	}
}

func (manager *SessionManager) Create(userID int64) (models.Result, error) {
	strUserID := fmt.Sprintf("%d", userID)

	session := &models.Session{
		UserID:     strUserID,
		UniqID:     generateUniqID(strUserID),
		Expiration: expiredCookiesTime,
	}
	if err := manager.sessionRepository.Set(session); err != nil {
		return models.Result{UserID: UnknownUser}, err
	}
	return models.Result{UserID: userID, UniqID: session.UniqID}, nil
}

func (manager *SessionManager) Delete(uniqID string) error {
	session := &models.Session{
		UniqID: uniqID,
	}
	return manager.sessionRepository.Del(session)
}

func (manager *SessionManager) Check(uniqID string) (models.Result, error) {
	userID, err := manager.sessionRepository.GetUserId(uniqID)
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
	value := append([]byte(userID), uuid.NewV4().Bytes()...)
	hash := md5.Sum(value)

	return hex.EncodeToString(hash[:])
}
