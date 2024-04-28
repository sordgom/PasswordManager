package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createPasswordRequest struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Hint     string `json:"hint"`
}

func (server *Server) createPassword(ctx *gin.Context) {
	var req createPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	vault := *server.VaultService.Vault
	// Load the vault from redis or check if it exists
	if vault.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Vault not found",
		})
		return
	}

	vault.NewPassword(req.Name, req.Url, req.Username, req.Password, req.Hint)
	server.VaultService.SaveVaultToRedis(&vault)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Password added successfully",
	})
}
