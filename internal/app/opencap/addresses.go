package opencap

import (
	"encoding/base64"
	"encoding/json"
	"github.com/opencap/opencap/internal/pkg/database"
	"github.com/opencap/opencap/pkg/messages"
	"github.com/valyala/fasthttp"
	"strconv"
)

func (s *Server) LookupAddress(ctx *fasthttp.RequestCtx) {
	var (
		domain    = ctx.UserValue("domain").(string)
		username  = ctx.UserValue("username").(string)
		typeIdStr = ctx.UserValue("type_id").(string)
	)

	typeId, err := strconv.ParseUint(typeIdStr, 10, 16)
	if err != nil {
		s.handleError(ctx, fasthttp.StatusBadRequest, "Invalid type_id")
		return
	}

	addr, err := s.db.GetAddress(domain, username, uint16(typeId))
	if err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Database query failed")
		return
	}
	if addr == nil {
		s.handleError(ctx, fasthttp.StatusNotFound, "No entry found")
		return
	}

	s.writeJSON(ctx, messages.LookupResponse{
		SubType:    addr.SubTypeId,
		Address:    base64.StdEncoding.EncodeToString(addr.AddressData),
		Extensions: map[string]interface{}{},
	})
}

func (s *Server) UpdateAddress(ctx *fasthttp.RequestCtx) {
	var (
		domain    = ctx.UserValue("domain").(string)
		username  = ctx.UserValue("username").(string)
		typeIdStr = ctx.UserValue("type_id").(string)
	)

	typeId, err := strconv.ParseUint(typeIdStr, 10, 16)
	if err != nil {
		s.handleError(ctx, fasthttp.StatusBadRequest, "Invalid type_id")
		return
	}

	var body messages.UpdateAddressRequest
	if err := json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		s.handleError(ctx, fasthttp.StatusBadRequest, "Failed to decode json")
		return
	}

	addressData, err := base64.StdEncoding.DecodeString(body.Address)
	if err != nil {
		s.handleError(ctx, fasthttp.StatusBadRequest, "Address must be valid base64")
		return
	}

	if err := s.db.SetAddress(domain, username, uint16(typeId), &database.Address{
		SubTypeId:   body.SubType,
		AddressData: addressData,
	}); err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Database write failed")
		return
	}
}

func (s *Server) DeleteAddress(ctx *fasthttp.RequestCtx) {
	var (
		domain   = ctx.UserValue("domain").(string)
		username = ctx.UserValue("username").(string)
	)

	if err := s.db.DeleteUser(domain, username); err != nil {
		s.handleError(ctx, fasthttp.StatusInternalServerError, "Database write failed")
		return
	}
}
