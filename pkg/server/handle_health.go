package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": "nook " + s.config.Version,
	})
}
