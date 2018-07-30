package opencap

import (
	"encoding/json"
	"github.com/opencap/opencap/pkg/messages"
	"github.com/valyala/fasthttp"
)

func (s *Server) AssociateDomain(ctx *fasthttp.RequestCtx) {
	var body messages.AssociateDomainRequest
	if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
		s.handleError(ctx, fasthttp.StatusBadRequest, "Invalid json")
		return
	}
}
