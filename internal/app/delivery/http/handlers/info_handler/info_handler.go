package info_handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"patreon/internal/app/delivery/http/handlers/base_handler"
	"patreon/internal/app/delivery/http/models"
	usecase_info "patreon/internal/app/usecase/info"
)

type InfoHandler struct {
	infoUsecase    usecase_info.Usecase
	base_handler.BaseHandler
}

func NewInfoHandler(log *logrus.Logger, ucInfo usecase_info.Usecase) *InfoHandler {
	h := &InfoHandler{
		BaseHandler:    *base_handler.NewBaseHandler(log),
		infoUsecase:    ucInfo,
	}

	h.AddMethod(http.MethodGet, h.GET)

	return h
}

// GET Info
// @Summary get info about creator category and type of post data
// @Description get info about creator category and type of post data
// @Produce json
// @Success 200 {object} models.ResponseInfo
// @Failure 500 {object} models.ErrResponse "can not do bd operation"
// @Failure 400 {object} models.ErrResponse "invalid parameters"
// @Router /info [GET]
func (s *InfoHandler) GET(w http.ResponseWriter, r *http.Request) {
	info, err := s.infoUsecase.Get()
	if err != nil {
		s.UsecaseError(w, r, err, codesByErrors)
		return
	}

	s.Log(r).Debug("get info")
	s.Respond(w, r, http.StatusOK, models.ToResponseInfo(*info))
}
