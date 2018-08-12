package opencap

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/opencap/opencap/internal/pkg/config"
	"github.com/opencap/opencap/internal/pkg/database"
	"github.com/opencap/opencap/pkg/messages"
	"github.com/valyala/fasthttp"
	"log"
)

type Server struct {
	db     database.Database
	config config.Config
}

func (s *Server) SetDatabase(db database.Database) {
	s.db = db
}

func (s *Server) SetConfig(config config.Config) {
	s.config = config
}

func (s *Server) Run() error {
	router := fasthttprouter.New()

	// CAP
	router.GET("/v1/domains/:domain/users/:username/types/:type_id", s.LookupAddress)

	// CAMP
	router.POST("/v1/domains/:domain/users", s.CreateUser)
	router.POST("/v1/auth", s.AuthenticateUser)
	router.DELETE("/v1/domains/:domain/users/:username", applyMiddleware(s.DeleteUser, s.requiresAuth, s.filterDomainAndUsernameMismatch))
	router.PUT("/v1/domains/:domain/users/:username/types/:type_id", applyMiddleware(s.UpdateAddress, s.requiresAuth, s.filterDomainAndUsernameMismatch))
	router.DELETE("/v1/domains/:domain/users/:username/types/:type_id", applyMiddleware(s.DeleteAddress, s.requiresAuth, s.filterDomainAndUsernameMismatch))

	// CAPP
	router.POST("/v1/domains", s.AssociateDomain)

	router.NotFound = s.handleNotImplemented

	log.Printf("Listening on %s:%d", s.config.Hostname(), s.config.Port())
	return fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Hostname(), s.config.Port()), router.Handler)
}

func (s *Server) writeJSON(ctx *fasthttp.RequestCtx, data interface{}) {
	enc := json.NewEncoder(ctx)

	if s.config.Debug() {
		enc.SetIndent("", "    ")
	}

	err := enc.Encode(data)
	if err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Internal server error")
	}
}

func (s *Server) handleError(ctx *fasthttp.RequestCtx, code int, msg string) {
	ctx.SetStatusCode(code)
	s.writeJSON(ctx, &messages.ErrorResponse{
		Code:    code,
		Message: msg,
	})
}

func (s *Server) handleNotImplemented(ctx *fasthttp.RequestCtx) {
	s.handleError(ctx, fasthttp.StatusNotImplemented, "Method not implemented")
}

func applyMiddleware(handler fasthttp.RequestHandler, middleware ...func(fasthttp.RequestHandler) fasthttp.RequestHandler) fasthttp.RequestHandler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}
