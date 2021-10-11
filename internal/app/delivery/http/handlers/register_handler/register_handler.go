package register_handler

import (
	"encoding/json"
	"io"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	usecaseModels "patreon/internal/app/models"
	"patreon/internal/app/sessions"
	"patreon/internal/app/sessions/middleware"
	useUser "patreon/internal/app/usecase/user"

	"github.com/sirupsen/logrus"
)

type RegisterHandler struct {
	sessionManager sessions.SessionsManager
	userUsecase    useUser.Usecase
	bh.BaseHandler
}

func NewRegisterHandler(log *logrus.Logger, sManager sessions.SessionsManager,
	ucUser useUser.Usecase) *RegisterHandler {
	h := &RegisterHandler{
		sessionManager: sManager,
		userUsecase:    ucUser,
		BaseHandler:    *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodPost, h.POST)
	h.AddMiddleware(middleware.NewSessionMiddleware(h.sessionManager, h.Log()).CheckNotAuthorized)
	return h
}

// Registration
// @Summary create new user
// @Description create new account and get cookies
// @Accept  json
// @Produce json
// @Param user body models.RequestRegistration true "Request body for user registration"
// @Success 201 {object} models.UserResponse "Create user successfully"
// @Failure 400 {object} models.BaseResponse "Invalid body"
// @Failure 409 {object} models.BaseResponse "User already exist"
// @Failure 422 {object} models.BaseResponse "Not valid body"
// @Failure 500 {object} models.BaseResponse "Creation error in base data"
// @Failure 418 "User are authorized"
// @Router /register [POST]
func (h *RegisterHandler) POST(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log().Error(err)
		}
	}(r.Body)
	req := &models.RequestRegistration{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil || len(req.Password) == 0 ||
		len(req.Nickname) == 0 || len(req.Login) == 0 {
		h.Log().Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}
	u := &usecaseModels.User{
		Login:    req.Login,
		Password: req.Password,
		Nickname: req.Nickname,
	}

	logUser, _ := json.Marshal(u)
	h.Log().Debug("get: ", string(logUser))

	id, err := h.userUsecase.Create(u)
	if err != nil {
		h.UsecaseError(w, r, err, codeByError)
		return
	}

	//checkUser, _ := h.dataStorage.Store().User().FindByLogin(u.Login)
	//if checkUser != nil {
	//	h.Log().Warn(handler_errors.UserAlreadyExist)
	//	h.Error(w, r, http.StatusConflict, handler_errors.UserAlreadyExist)
	//	return
	//}
	//if err := u.Validate(); err != nil {
	//	h.Log().Warnf("Not valid login or password %s", err)
	//	h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidBody)
	//	return
	//}
	//
	//if err := u.Encrypt(); err != nil {
	//	h.Log().Errorf("Error prepare user info %s", err)
	//	h.Error(w, r, http.StatusInternalServerError, handler_errors.ErrorPrepareUser)
	//	return
	//}
	//
	//if err := h.dataStorage.Store().User().Create(u); err != nil {
	//	h.Log().Errorf("Error create user in bd %s", err)
	//	h.Error(w, r, http.StatusInternalServerError, handler_errors.ErrorCreateUser)
	//	return
	//}
	u.MakePrivateDate()
	h.Respond(w, r, http.StatusOK, id)
}
