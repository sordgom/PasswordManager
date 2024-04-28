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

type createPasswordVaultName struct {
	Name string `form:"vault_name" binding:"required"`
}

func (server *Server) createPassword(ctx *gin.Context) {
	var req createPasswordRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var params createPasswordVaultName
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Load the vault from redis or check if it exists
	vault, err := server.VaultService.LoadVaultFromRedis(ctx, params.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Vault not found",
		})
		return
	}

	vault.NewPassword(req.Name, req.Url, req.Username, req.Password, req.Hint)
	server.VaultService.SaveVaultToRedis(ctx, vault)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Password added successfully",
	})
}

func (server *Server) getPasswords(ctx *gin.Context) {
	var params createPasswordVaultName
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	vault, err := server.VaultService.LoadVaultFromRedis(ctx, params.Name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Vault not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, vault.Passwords)
}
