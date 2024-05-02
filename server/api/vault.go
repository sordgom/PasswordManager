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

	if req.Name == "" || req.MasterPassword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Name and master password are required",
		})
		return
	}

	vault := model.Vault{
		Name:           req.Name,
		MasterPassword: req.MasterPassword,
	}

	// Check if the vault already exists
	_, err := server.VaultService.LoadVaultFromRedis(ctx, vault.Name)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Vault already exists",
		})
		return
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

type verifyMasterPasswordRequest struct {
	MasterPassword string `json:"master_password"`
	VaultName      string `json:"vault_name"`
}

func (server *Server) verifyMasterPassword(ctx *gin.Context) {
	var req verifyMasterPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.VaultName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Vault name is required",
		})
		return
	}

	// Verify the master password
	isVerified := server.VaultService.VerifyMasterPassword(ctx, req.VaultName, req.MasterPassword)

	if isVerified {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Master password is verified",
		})
		return
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Master password is not verified",
		})
		return
	}
}
