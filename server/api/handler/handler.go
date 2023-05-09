package handler

import (
	apimiddleware "github.com/Xrefullx/YanDip/server/api/middleware"
	"github.com/Xrefullx/YanDip/server/services/auth"
	"github.com/Xrefullx/YanDip/server/services/secret"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

type Handler struct {
	svcAuth   auth.Authenticator
	svcSecret secret.SecretManager
	jwtAuth   *jwtauth.JWTAuth
}

// NewHandler Return new handler
func NewHandler(auth auth.Authenticator, secret secret.SecretManager, jwtAuth *jwtauth.JWTAuth) (*Handler, error) {

	return &Handler{
		svcAuth:   auth,
		svcSecret: secret,
		jwtAuth:   jwtAuth,
	}, nil
}

func GetRouter(handler *Handler) *chi.Mux {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	r := chi.NewRouter()
	r.Use(middleware.Compress(5))

	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/api/user/register", handler.Register)
		r.Post("/api/user/login", handler.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(apimiddleware.MiddlewareAuth)

		r.Get("/api/sync", handler.SyncList)
		r.Get("/api/ping", handler.Ping)

		// Secret processing
		r.Put("/api/secret", handler.SecretUpload)
		r.Get("/api/secret", handler.SecretGet)
		r.Delete("/api/secret", handler.SecretDelete)

	})

	return r
}

// Ping returns 200
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	h.writeJSONResponse(w, http.StatusOK, "ok")
}
