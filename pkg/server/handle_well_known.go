package server

import "github.com/gin-gonic/gin"

func (s *Server) handleWellKnown(c *gin.Context) {
	c.JSON(200, gin.H{
		"@context": []string{
			"https://www.w3.org/ns/did/v1",
		},
		"id": s.config.Did,
		"service": []gin.H{
			{
				"id":              "#atproto_pds",
				"type":            "AtprotoPersonalDataServer",
				"serviceEndpoint": "https://" + s.config.Hostname,
			},
		},
	})
}
