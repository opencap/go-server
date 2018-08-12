package opencap

import (
	"encoding/json"
	"github.com/opencap/opencap/pkg/messages"
	"github.com/opencap/opencap/pkg/resolver"
	"github.com/valyala/fasthttp"
)

func (s *Server) AssociateDomain(ctx *fasthttp.RequestCtx) {
	var body messages.AssociateDomainRequest
	if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
		s.handleError(ctx, fasthttp.StatusBadRequest, "Invalid json")
		return
	}

	res, err := resolver.Resolve(body.Domain)
	if err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Failed to resolve "+body.Domain)
		return
	}

	if res.PublicKey != nil {
		if err := s.db.SetPublicKey(body.Domain, res.PublicKey); err != nil {
			s.handleError(ctx, fasthttp.StatusInternalServerError, "Failed to set/update public key")
			return
		}
	} else {
		if err := s.db.DeletePublicKey(body.Domain); err != nil {
			s.handleError(ctx, fasthttp.StatusInternalServerError, "Failed to delete public key")
			return
		}
	}
}
