package search_creators_handler

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/handlers/handler_errors"
	"patreon/internal/app/delivery/http/models"
	"patreon/internal/app/models"
	usecase_creator "patreon/internal/app/usecase/creator"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
)

type SearchCreatorsHandler struct {
	sessionClient  session_client.AuthCheckerClient
	creatorUsecase usecase_creator.Usecase
	bh.BaseHandler
}

func NewCreatorHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient,
	ucCreator usecase_creator.Usecase) *SearchCreatorsHandler {
	h := &SearchCreatorsHandler{
		BaseHandler:    *bh.NewBaseHandler(log),
		creatorUsecase: ucCreator,
		sessionClient:  sManager,
	}
	h.AddMethod(http.MethodGet, h.GET)
	return h
}

// GET Creators
// @Summary get list of Creators
// @Description get list of creators which register on service
// @Produce json
// @tags creators
// @Param page query uint64 true "start page number of creators mutually exclusive with offset"
// @Param offset query uint64 true "start number of creators mutually exclusive with page"
// @Param limit query uint64 true "creators to return"
// @Param search_string query string true "search query"
// @Param category query string false "if need filter category, may be several param in one query"
// @Success 201 {array} http_models.ResponseCreators
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation"
// @Failure 400 {object} http_models.ErrResponse "invalid parameters in query"
// @Router /creators/search [GET]
func (h *SearchCreatorsHandler) GET(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := h.GetPaginationFromQuery(w, r)
	if !ok {
		return
	}

	var searchString string
	if searchString = r.URL.Query().Get("search_string"); searchString == "" {
		h.Error(w, r, http.StatusBadRequest, handler_errors.InvalidQueries)
		return
	}
	sanitizer := bluemonday.NewPolicy()
	searchString = sanitizer.Sanitize(searchString)

	categories := r.URL.Query()["category"]
	for i, cat := range categories {
		categories[i] = sanitizer.Sanitize(cat)
	}

	creators, err := h.creatorUsecase.SearchCreators(&models.Pagination{
		Limit: limit, Offset: offset}, searchString, categories...)
	if err != nil {
		h.UsecaseError(w, r, err, codesByErrorsGET)
		return
	}

	respondCreators := http_models.ToResponseCreators(creators)

	h.Log(r).Debugf("get creators %v with search query %s and categories %v",
		respondCreators, searchString, categories)
	h.Respond(w, r, http.StatusOK, respondCreators)
}
