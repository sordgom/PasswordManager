package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) createVault(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "createVault",
	})
}
