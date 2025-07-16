package helpers

import "github.com/gin-gonic/gin"

func InputError(c *gin.Context, custom *string) {
	msg := "InvalidRequest"
	if custom != nil {
		msg = *custom
	}
	genericError(c, 400, msg)
}

func ServerError(c *gin.Context, suffix *string) {
	msg := "Internal server error"
	if suffix != nil {
		msg += ", " + *suffix
	}
	genericError(c, 500, msg)
}

func genericError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"error": msg,
	})
}
