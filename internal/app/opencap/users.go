package opencap

import (
	"encoding/json"
	"github.com/opencap/opencap/pkg/messages"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) CreateUser(ctx *fasthttp.RequestCtx) {
	var body messages.CreateUserRequest
	if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
		s.handleError(ctx, fasthttp.StatusBadRequest, "Invalid json")
		return
	}

	var (
		domain   = ctx.UserValue("domain").(string)
		username = body.Username
		password = body.Password
	)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Hashing password failed")
		return
	}

	if err := s.db.CreateUser(domain, username, string(hash)); err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Database write failed")
		return
	}
}

func (s *Server) DeleteUser(ctx *fasthttp.RequestCtx) {
	var (
		claims   = ctx.UserValue("jwt").(*Claims)
		username = ctx.UserValue("username").(string)
		domain   = ctx.UserValue("domain").(string)
	)

	if claims.Domain != domain {
		s.handleError(ctx, fasthttp.StatusUnauthorized, "Domains do not match")
		return
	}

	if claims.Username != username {
		s.handleError(ctx, fasthttp.StatusUnauthorized, "Usernames do not match")
		return
	}

	if err := s.db.DeleteUser(domain, username); err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Database write failed")
		return
	}
}
