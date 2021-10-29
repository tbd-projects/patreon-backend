package posts_handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"patreon/internal/app"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	db_models "patreon/internal/app/models"
	"patreon/internal/app/sessions"
	sessionMid "patreon/internal/app/sessions/middleware"
	usePosts "patreon/internal/app/usecase/posts"
)

type PostsHandler struct {
	postsUsecase usePosts.Usecase
	bh.BaseHandler
}

func NewPostsHandler(log *logrus.Logger, router *mux.Router, cors *app.CorsConfig,
	ucPosts usePosts.Usecase, manager sessions.SessionsManager) *PostsHandler {
	h := &PostsHandler{
		BaseHandler:  *bh.NewBaseHandler(log, router, cors),
		postsUsecase: ucPosts,
	}
	sessionMiddleware := sessionMid.NewSessionMiddleware(manager, log)
	h.AddMethod(http.MethodGet, h.GET)
	h.AddMethod(http.MethodPost, h.POST, sessionMiddleware.CheckFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc)
	return h
}

// GETRedirect Posts
// @Summary redirect to get post with default query
// @Description redirect to get post with default query
// @Produce json
// @Success 308
// @Router /creators/{:creator_id}/posts [GET]
func (h *PostsHandler) redirect(w http.ResponseWriter, r *http.Request) {
	redirectUrl := fmt.Sprintf("%s?page=1&limit=%d", r.RequestURI, usePosts.BaseLimit)
	h.Log(r).Debugf("redirect to url: %s, with offest 0 and limit %d", redirectUrl, usePosts.BaseLimit)

	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}

func (h *PostsHandler) baseGet(w http.ResponseWriter, r *http.Request, pag *db_models.Pagination) {

}

// GET Posts
// @Summary get list of posts of some creator
// @Description get list of posts which belongs the creator with limit and offset in query
// @Produce json
// @Success 201 {array} models.ResponsePost
// @Param page query uint64 true "start page number of posts mutually exclusive with offset"
// @Param offset query uint64 true "start number of posts mutually exclusive with page"
// @Param limit query uint64 true "posts to return"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 400 {object} models.ErrResponse "invalid parameters in query"
// @Failure 500 {object} models.ErrResponse "can not get info from context"
// @Router /creators/{:creator_id}/posts [GET]
func (h *PostsHandler) GET(w http.ResponseWriter, r *http.Request) {
	var limit, offset, page int64
	var ok bool

	limit, ok = h.GetInt64FromQueries(w, r, "limit")
	if !ok {
		if limit == bh.EmptyQuery {
			h.redirect(w, r)
		}
		return
	}

	offset, ok = h.GetInt64FromQueries(w, r, "offset")
	if !ok {
		if offset != bh.EmptyQuery {
			return
		}
		page, ok = h.GetInt64FromQueries(w, r, "page")
		if !ok {
			if offset == bh.EmptyQuery {
				h.redirect(w, r)
			}
			return
		}
		if page <= 0 {
			page = 1
		}
		offset = (page - 1) * limit
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	var creatorId, userId int64

	creatorId, ok = h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	userId, ok = r.Context().Value("user_id").(int64)
	if !ok {
		userId = usePosts.EmptyUser
	}

	posts, err := h.postsUsecase.GetPosts(creatorId, userId, &db_models.Pagination{Limit: limit, Offset: offset})
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondPosts := make([]models.ResponsePost, len(posts))
	for i, ps := range posts {
		respondPosts[i] = models.ToResponsePost(ps)
	}

	h.Log(r).Debugf("get posts %v", respondPosts)
	h.Respond(w, r, http.StatusOK, respondPosts)
}

// POST Create Posts
// @Summary create posts
// @Description create posts to creator with id from path
// @Param user body models.RequestPosts true "Request body for posts"
// @Produce json
// @Success 201 {object} models.IdResponse "id posts"
// @Failure 422 {object} models.ErrResponse "invalid body in request"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Failure 422 {object} models.ErrResponse "empty title"
// @Failure 422 {object} models.ErrResponse "this creator id not know"
// @Failure 422 {object} models.ErrResponse "this awards id not know"
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 500 {object} models.ErrResponse "server error"
// @Failure 500 {object} models.ErrResponse "can not get info from context"
// @Failure 403 {object} models.ErrResponse "for this user forbidden change creator"
// @Failure 401 "User are not authorized"
// @Router /creators/{:creator_id}/posts [POST]
func (h *PostsHandler) POST(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Log(r).Error(err)
		}
	}(r.Body)

	idInt, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	req := &models.RequestPosts{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	aw := &db_models.CreatePost{
		Title:       req.Title,
		Description: req.Description,
		Awards:      req.AwardsId,
		CreatorId:   idInt,
	}

	postsId, err := h.postsUsecase.Create(aw)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}

	h.Respond(w, r, http.StatusCreated, &models.IdResponse{ID: postsId})
}
