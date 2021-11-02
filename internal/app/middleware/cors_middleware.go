package middleware

import (
	"net/http"
	"patreon/internal/app"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type CorsMiddleware struct {
	router *mux.Router
	config *app.CorsConfig
}

func NewCorsMiddleware(config *app.CorsConfig, router *mux.Router) CorsMiddleware {
	return CorsMiddleware{
		router: router,
		config: config,
	}
}
func (mw *CorsMiddleware) SetCors(handler http.Handler) http.Handler {
	return gh.CORS(
			gh.AllowedOrigins(mw.config.Urls),
			gh.AllowedHeaders(mw.config.Headers),
			gh.AllowCredentials(),
			gh.AllowedMethods(mw.config.Methods),
		) (handler)
}
