package server

import (
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/gin-gonic/gin"
	"github.com/hbjydev/nook/internal/helpers"
	"github.com/hbjydev/nook/pkg/identity"
)

func (s *Server) handleIdentityResolveHandle(c *gin.Context) {
	type Resp struct {
		Did string `json:"did"`
	}

	handle := c.Query("handle")
	if handle == "" {
		helpers.InputError(c, to.StringPtr("Handle must be supplied in request."))
		return
	}

	parsed, err := syntax.ParseHandle(handle)
	if err != nil {
		helpers.InputError(c, to.StringPtr("Invalid handle."))
		return
	}

	did, err := identity.ResolveHandle(c.Request.Context(), s.http, parsed.String())
	if err != nil {
		s.logger.Error("error resolving handle", "error", err)
		helpers.ServerError(c, nil)
		return
	}

	c.JSON(200, Resp{
		Did: did,
	})
}
