package opencap

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/opencap/opencap/pkg/messages"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

const defTokenLifetime = 600 // 10 minutes

type Claims struct {
	Username string `json:"username"`
	Domain   string `json:"domain"`
	jwt.StandardClaims
}

func (s *Server) AuthenticateUser(ctx *fasthttp.RequestCtx) {
	var body messages.LoginRequest
	if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
		s.handleError(ctx, fasthttp.StatusBadRequest, "Invalid json")
		return
	}

	hashedPassword, err := s.db.GetUserPassword(body.Domain, body.Username)
	if err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Database query failed")
		return
	}
	if len(hashedPassword) == 0 {
		s.handleError(ctx, fasthttp.StatusNotFound, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body.Password)); err != nil {
		s.handleError(ctx, fasthttp.StatusUnauthorized, "Invalid password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username: body.Username,
		Domain:   body.Domain,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + defTokenLifetime,
		},
	})

	tokenString, err := token.SignedString("")
	if err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Signing JWT failed")
		return
	}

	s.writeJSON(ctx, messages.LoginResponse{
		JWT: tokenString,
	})
}

func (s *Server) requiresAuth(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	const prefix = "Bearer "

	return func(ctx *fasthttp.RequestCtx) {
		auth := string(ctx.Request.Header.Peek("Authorization"))

		if !strings.HasPrefix(auth, prefix) {
			s.handleError(ctx, fasthttp.StatusBadRequest, "Bad auth header")
			return
		}

		auth = auth[len(prefix):]

		token, err := jwt.ParseWithClaims(auth, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return token, nil
		})
		if err != nil {
			s.handleError(ctx, fasthttp.StatusInternalServerError, "Parsing token failed")
			return
		}

		claims := token.Claims.(*Claims)

		if err := claims.Valid(); err != nil {
			s.handleError(ctx, fasthttp.StatusUnauthorized, "Invalid JWT claims: "+err.Error())
			return
		}

		ctx.SetUserValue("jwt", claims)

		handler(ctx)
	}
}

func (s *Server) filterDomainAndUsernameMismatch(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
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

		handler(ctx)
	}
}
