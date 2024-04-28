package api

import (
	"net/http"

	"github.com/sordgom/PasswordManager/server/model"

	"github.com/gin-gonic/gin"
)

type createVaultRequest struct {
	Name           string `json:"name"`
	MasterPassword string `json:"master_password"`
}

func (server *Server) createVault(ctx *gin.Context) {
	var req createVaultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	vault := model.Vault{
		Name:           req.Name,
		MasterPassword: req.MasterPassword,
	}

	// Save the vault using the service
	if err := server.VaultService.SaveVaultToRedis(ctx, &vault); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Vault created successfully",
	})
}
