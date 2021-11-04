package upd_text_data_handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	models_db "patreon/internal/app/models"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	usePosts "patreon/internal/app/usecase/posts"
	usePostsData "patreon/internal/app/usecase/posts_data"

	"github.com/sirupsen/logrus"
)

type PostsDataUpdateTextHandler struct {
	postsDataUsecase usePostsData.Usecase
	bh.BaseHandler
}

func NewPostsDataUpdateTextHandler(log *logrus.Logger,
	ucPostsData usePostsData.Usecase, ucPosts usePosts.Usecase,
	manager sessions.SessionsManager) *PostsDataUpdateTextHandler {
	h := &PostsDataUpdateTextHandler{
		BaseHandler:      *bh.NewBaseHandler(log),
		postsDataUsecase: ucPostsData,
	}
	sessionMiddleware := sessionMid.NewSessionMiddleware(manager, log)
	h.AddMiddleware(sessionMiddleware.Check, middleware.NewCreatorsMiddleware(log).CheckAllowUser,
		middleware.NewPostsMiddleware(log, ucPosts).CheckCorrectPost, sessionMiddleware.AddUserId)
	h.AddMethod(http.MethodPut, h.PUT)
	return h
}

// PUT update text to post
// @Summary update text to post
// @Accept  json
// @Param user body models.RequestText true "Request body for text"
// @Success 201 {object} models.IdResponse "id posts_data"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error"
// @Failure 422 {object} models.ErrResponse "invalid data type"
// @Failure 422 {object} models.ErrResponse "this post id not know"
// @Failure 404 {object} models.ErrResponse "post data with this id not found"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Router /creators/{:creator_id}/posts/{:post_id}/{:data_id}/update/text [PUT]
func (h *PostsDataUpdateTextHandler) PUT(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)

	var dataId int64
	var ok bool

	if dataId, ok = h.GetInt64FromParam(w, r, "data_id"); !ok {
		return
	}

	if len(mux.Vars(r)) > 3 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	req := &models.RequestText{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	err := h.postsDataUsecase.UpdateText(&models_db.PostData{ID: dataId, Data: req.Text})
	if err != nil {
		h.UsecaseError(w, r, err, codeByErrorPUT)
		return
	}

	w.WriteHeader(http.StatusOK)
}
