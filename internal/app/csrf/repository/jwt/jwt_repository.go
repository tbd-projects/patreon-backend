package repository_jwt

import (
	"patreon/internal/app"
	"patreon/internal/app/csrf/models"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/golang-jwt/jwt"
)

var (
	secretToken = uuid.NewV4()
)

type JwtRepository struct {
	Secret []byte
}

type jwtCsrfClaims struct {
	jwt.StandardClaims
	UserId    int64  `json:"user_id"`
	SessionId string `json:"session_id"`
}

func NewJwtToken() *JwtRepository {
	return &JwtRepository{Secret: secretToken.Bytes()}
}

func (tk *JwtRepository) parseClaims(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, IncorrectTokenSigningMethod
	}
	return tk.Secret, nil
}
func (tk *JwtRepository) Check(sources models.TokenSources, tokenString models.Token) error {
	claims := &jwtCsrfClaims{}
	if _, err := jwt.ParseWithClaims(string(tokenString), claims, tk.parseClaims); err != nil {
		return &app.GeneralError{
			Err:         ParseClaimsError,
			ExternalErr: err,
		}
	}
	if err := claims.Valid(); err != nil {
		return &app.GeneralError{
			Err:         TokenExpired,
			ExternalErr: err,
		}
	}
	if claims.UserId != sources.UserId || claims.SessionId != sources.SessionId {
		return BadToken
	}
	return nil
}
func (tk *JwtRepository) Create(sources models.TokenSources) (models.Token, error) {
	data := jwtCsrfClaims{
		SessionId: sources.SessionId,
		UserId:    sources.UserId,
		StandardClaims: jwt.StandardClaims{
			//ExpiresAt: time.Now().Add(expiredCsrfTime).Unix(),
			ExpiresAt: sources.ExpiredTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	res, err := token.SignedString(tk.Secret)
	if err != nil {
		return "", &app.GeneralError{
			Err:         ErrorSignedToken,
			ExternalErr: err,
		}
	}
	return models.Token(res), nil
}
