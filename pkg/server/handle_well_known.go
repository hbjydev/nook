package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hbjydev/nook/pkg/identity"
)

func (s *Server) handleWellKnown(c *gin.Context) {
	c.JSON(200, identity.DidDoc{
		Context: []string{"https://www.w3.org/ns/did/v1"},
		Id: s.config.Did,
		Service: []identity.DidDocService{
			{
				Id: "#atproto_pds",
				Type: "AtprotoPersonalDataServer",
				ServiceEndpoint: "https://" + s.config.Hostname,
			},
		},
	})
}

func (s *Server) handleOauthProtectedResource(c *gin.Context) {
	c.JSON(200, map[string]any{
		"resource": "https://" + s.config.Hostname,
		"authorization_servers": []string{
			"https://" + s.config.Hostname,
		},
		"scopes_supported":         []string{},
		"bearer_methods_supported": []string{"header"},
		"resource_documentation":   "https://atproto.com",
	})
}
