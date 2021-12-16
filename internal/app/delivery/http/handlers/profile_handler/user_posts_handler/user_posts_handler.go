package user_posts_handler

import (
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	app_models "patreon/internal/app/models"
	usecase_posts "patreon/internal/app/usecase/posts"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	session_middleware "patreon/internal/microservices/auth/sessions/middleware"

	"github.com/sirupsen/logrus"
)

type PostsHandler struct {
	sessionClient session_client.AuthCheckerClient
	postsUsecase  usecase_posts.Usecase
	bh.BaseHandler
}

func NewPostsHandler(log *logrus.Logger, sClient session_client.AuthCheckerClient,
	ucPosts usecase_posts.Usecase) *PostsHandler {
	h := &PostsHandler{
		sessionClient: sClient,
		postsUsecase:  ucPosts,
		BaseHandler:   *bh.NewBaseHandler(log),
	}
	h.AddMethod(http.MethodGet, h.GET,
		session_middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc)
	return h
}

// GET AvailablePosts
// @Summary get all user available posts
// @tags user
// @Description get user available posts
// @Produce json
// @Param page query uint64 true "start page number of posts mutually exclusive with offset"
// @Param offset query uint64 true "start number of posts mutually exclusive with page"
// @Param limit query uint64 true "posts to return"
// @Success 200 {object} http_models.ResponseAvailablePosts "Successfully get user available posts"
// @Success 204  "No available posts"
// @Failure 500 {object} http_models.ErrResponse "serverError"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters", "invalid parameters in query"
// @Failure 401 "user are not authorized"
// @Router /user/posts [GET]
func (h *PostsHandler) GET(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := h.GetPaginationFromQuery(w, r)
	if !ok {
		return
	}
	userID := r.Context().Value("user_id")
	if userID == nil {
		h.Log(r).Error("can not get user_id from context")
		h.Error(w, r, http.StatusInternalServerError, handler_errors.InternalError)
		return
	}
	posts, err := h.postsUsecase.GetAvailablePosts(userID.(int64), &app_models.Pagination{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		h.UsecaseError(w, r, err, codesByErrors)
		return
	}

	if len(posts) == 0 {
		h.Log(r).Errorf("no available posts for user with id = %d", userID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	res := http_models.ToResponseAvailablePosts(posts)
	h.Log(r).Debugf("get available posts %v", posts)
	h.Respond(w, r, http.StatusOK, res)
}
