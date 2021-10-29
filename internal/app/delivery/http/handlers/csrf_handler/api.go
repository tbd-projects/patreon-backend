package csrf_handler

import (
	"net/http"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"

	"github.com/sirupsen/logrus"
)

var codeByErrors = base_handler.CodeMap{
	repository_jwt.ErrorSignedToken: {http.StatusInternalServerError, handler_errors.InternalError, logrus.ErrorLevel},
}
