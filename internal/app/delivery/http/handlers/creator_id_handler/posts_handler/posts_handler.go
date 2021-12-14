package posts_handler

import (
	"net/http"
	csrf_middleware "patreon/internal/app/csrf/middleware"
	repository_jwt "patreon/internal/app/csrf/repository/jwt"
	usecase_csrf "patreon/internal/app/csrf/usecase"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/middleware"
	db_models "patreon/internal/app/models"
	usePosts "patreon/internal/app/usecase/posts"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"
	"strconv"

	"github.com/microcosm-cc/bluemonday"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PostsHandler struct {
	postsUsecase usePosts.Usecase
	bh.BaseHandler
}

func NewPostsHandler(log *logrus.Logger,
	ucPosts usePosts.Usecase, sClient session_client.AuthCheckerClient) *PostsHandler {
	h := &PostsHandler{
		BaseHandler:  *bh.NewBaseHandler(log),
		postsUsecase: ucPosts,
	}
	sessionMiddleware := session_middleware.NewSessionMiddleware(sClient, log)

	h.AddMethod(http.MethodGet, h.GET)

	h.AddMethod(http.MethodPost, h.POST, sessionMiddleware.CheckFunc,
		middleware.NewCreatorsMiddleware(log).CheckAllowUserFunc,
		csrf_middleware.NewCsrfMiddleware(log,
			usecase_csrf.NewCsrfUsecase(repository_jwt.NewJwtRepository())).CheckCsrfTokenFunc)

	return h
}

// GET Posts
// @Summary get list of posts of some creator
// @tags posts
// @Description get list of posts which belongs the creator with limit and offset in query
// @Produce json
// @Success 201 {array} http_models.ResponsePost
// @Param page query uint64 true "start page number of posts mutually exclusive with offset"
// @Param offset query uint64 true "start number of posts mutually exclusive with page"
// @Param limit query uint64 true "posts to return"
// @Param with-draft query bool false "if need add draft posts"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters", "invalid parameters in query"
// @Router /creators/{:creator_id}/posts [GET]
func (h *PostsHandler) GET(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := h.GetPaginationFromQuery(w, r)
	if !ok {
		return
	}

	var err error
	var withDraft bool
	if res := r.URL.Query().Get("with-draft"); res == "" {
		withDraft = false
	} else if withDraft, err = strconv.ParseBool(res); err != nil {
		withDraft = true
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

	posts, err := h.postsUsecase.GetPosts(creatorId, userId,
		&db_models.Pagination{Limit: limit, Offset: offset}, withDraft)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondPosts := make([]http_models.ResponsePost, len(posts))
	for i, ps := range posts {
		respondPosts[i] = http_models.ToResponsePost(ps)
	}

	h.Log(r).Debugf("get posts %v", respondPosts)
	h.Respond(w, r, http.StatusOK, respondPosts)
}

// POST Create Posts
// @Summary create posts
// @tags posts
// @Description create posts to creator with id from path
// @Param post body http_models.RequestPosts true "Request body for posts"
// @Produce json
// @Success 201 {object} http_models.IdResponse "id posts"
// @Failure 400 {object} http_models.ErrResponse "invalid body in request"
// @Failure 422 {object} http_models.ErrResponse "this creator id not know", "this awards id not know", "empty title", "invalid parameters"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 403 {object} http_models.ErrResponse "for this user forbidden change creator", "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /creators/{:creator_id}/posts [POST]
func (h *PostsHandler) POST(w http.ResponseWriter, r *http.Request) {
	req := &http_models.RequestPosts{}

	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not parse request %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler_errors.InvalidBody)
		return
	}

	idInt, ok := h.GetInt64FromParam(w, r, "creator_id")
	if !ok {
		return
	}

	if len(mux.Vars(r)) > 1 {
		h.Log(r).Warnf("Too many parametres %v", mux.Vars(r))
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidParameters)
		return
	}

	aw := &db_models.CreatePost{
		Title:       req.Title,
		Description: req.Description,
		Awards:      req.AwardsId,
		CreatorId:   idInt,
		IsDraft:     req.IsDraft,
	}

	postsId, err := h.postsUsecase.Create(h.Log(r), aw)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsPOST)
		return
	}

	h.Respond(w, r, http.StatusCreated, &http_models.IdResponse{ID: postsId})
}
