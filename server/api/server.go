package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/Xrefullx/YanDip/server/api/handler"
	"github.com/Xrefullx/YanDip/server/pkg"
	"github.com/Xrefullx/YanDip/server/services/auth"
	"github.com/Xrefullx/YanDip/server/services/secret"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
)

type Server struct {
	httpServer http.Server
	cfg        *pkg.Config
}

func NewServer(cfg *pkg.Config, a auth.Authenticator, secret secret.SecretManager, jwtAuth *jwtauth.JWTAuth) (*Server, error) {

	h, err := handler.NewHandler(a, secret, jwtAuth)
	if err != nil {
		return nil, fmt.Errorf("ошибка запуска server:%w", err)
	}

	return &Server{
		httpServer: http.Server{
			Addr:    cfg.ServerPort,
			Handler: handler.GetRouter(h),
		},
		cfg: cfg,
	}, nil

}

// Run Start server
func (s *Server) Run() error {
	log.Printf("starting HTTP server on %v", s.cfg.ServerPort)

	if s.cfg.EnableHTTPS {
		certPath, keyPath, err := pkg.GetCertX509Files()
		if err != nil {
			return fmt.Errorf("error serve ssl:%w", err)
		}
		return handleServerCloseErr(s.httpServer.ListenAndServeTLS(certPath, keyPath))
	}

	return handleServerCloseErr(s.httpServer.ListenAndServe())
}

// returns error if error is not http.ErrServerClosed
func handleServerCloseErr(err error) error {
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server closed with: %w", err)
	}

	return nil
}

// Shutdown server
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.ListenAndServe(); err == http.ErrServerClosed {
		return errors.New("http server not runned")
	}

	return s.httpServer.Shutdown(ctx)
}
