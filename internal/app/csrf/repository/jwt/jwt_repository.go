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

func NewJwtRepository() *JwtRepository {
	return &JwtRepository{Secret: secretToken.Bytes()}
}

func (tk *JwtRepository) parseClaims(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok || method.Alg() != "HS256" {
		return nil, IncorrectTokenSigningMethod
	}
	return tk.Secret, nil
}

// Check Errors:
// 		repository_jwt.BadToken
// 		app.GeneralError with Error
// 			repository_jwt.ParseClaimsError
// 			repository_jwt.TokenExpired
func (tk *JwtRepository) Check(sources csrf_models.TokenSources, tokenString csrf_models.Token) error {
	claims := &jwtCsrfClaims{}
	token, err := jwt.ParseWithClaims(string(tokenString), claims, tk.parseClaims)
	if err != nil || !token.Valid {
		retErr := &app.GeneralError{ExternalErr: err}

		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			retErr.Err = TokenExpired
		case jwt.ValidationErrorUnverifiable:
			retErr.ExternalErr = IncorrectTokenSigningMethod
			retErr.Err = ParseClaimsError
		default:
			retErr.Err = ParseClaimsError
		}
		return retErr
	}

	if claims.UserId != sources.UserId || claims.SessionId != sources.SessionId {
		return BadToken
	}
	return nil
}

// Create Errors:
// 		app.GeneralError with Error
// 			repository_jwt.ErrorSignedToken
func (tk *JwtRepository) Create(sources csrf_models.TokenSources) (csrf_models.Token, error) {
	data := jwtCsrfClaims{
		SessionId: sources.SessionId,
		UserId:    sources.UserId,
		StandardClaims: jwt.StandardClaims{
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
	return csrf_models.Token(res), nil
}
