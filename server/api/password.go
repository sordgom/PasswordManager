package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) createPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "createPassword",
	})
}
