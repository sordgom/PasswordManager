package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sordgom/PasswordManager/server/model"
)

type createPasswordRequest struct {
	Name     string `json:"name" binding:"required"`
	Url      string `json:"url"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hint     string `json:"hint"`
}

type createPasswordVaultName struct {
	Name string `form:"vault_name" binding:"required"`
}

func (server *Server) createPassword(ctx *gin.Context) {
	var req createPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

type getPasswordsResponse struct {
	Name string `json:"name"`
	Hint string `json:"hint"`
}

// Convert the password struct to getPasswordsResponse struct
// Return the list of passwords in the vault
func ToPasswordsResponse(passwords []model.Password) []getPasswordsResponse {
	var result []getPasswordsResponse
	for _, password := range passwords {
		newPassword := getPasswordsResponse{
			Name: password.Name,
			Hint: password.Hint,
		}
		result = append(result, newPassword)
	}
	return result
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

	passwords := ToPasswordsResponse(vault.Passwords)

	ctx.JSON(http.StatusOK, passwords)
}

type getPasswordRequest struct {
	VaultName    string `form:"vault_name" binding:"required"`
	PasswordName string `form:"password_name" binding:"required"`
}

type getPasswordResponse struct {
	Name     string `json:"name"`
	Hint     string `json:"hint"`
	Password string `json:"password"`
}

func (server *Server) getPassword(ctx *gin.Context) {
	var params getPasswordRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// if !server.VaultService.VerifyMasterPassword(ctx, params.VaultName, params.MasterPassword) {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Invalid master password",
	// 	})
	// 	return
	// }

	vault, err := server.VaultService.LoadVaultFromRedis(ctx, params.VaultName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Vault not found",
		})
		return
	}

	password, err := vault.GetPassword(params.PasswordName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Password not found",
		})
		return
	}

	passwordResponse := getPasswordResponse{
		Name:     password.Name,
		Hint:     password.Hint,
		Password: vault.ReadPassword(&password),
	}

	ctx.JSON(http.StatusOK, passwordResponse)
}

type getPasswordByUrlRequest struct {
	VaultName string `form:"vault_name" binding:"required"`
	Url       string `form:"url"`
}

func (server *Server) getPasswordByUrl(ctx *gin.Context) {
	var params getPasswordByUrlRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	vault, err := server.VaultService.LoadVaultFromRedis(ctx, params.VaultName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Vault not found",
		})
		return
	}

	password, err := vault.GetPasswordByUrl(params.Url)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	passwordResponse := getPasswordResponse{
		Name:     password.Name,
		Hint:     password.Hint,
		Password: vault.ReadPassword(&password),
	}

	ctx.JSON(http.StatusOK, passwordResponse)
}

type updatePasswordRequest struct {
	VaultName       string `json:"vault_name"` //not sure if this should be in query or body.
	Name            string `json:"name" binding:"required"`
	MasterPassword  string `json:"master_password" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

func (server *Server) updatePassword(ctx *gin.Context) {
	var req updatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	vault, err := server.VaultService.LoadVaultFromRedis(ctx, req.VaultName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Vault not found",
		})
		return
	}

	if !server.VaultService.VerifyMasterPassword(ctx, req.VaultName, req.MasterPassword) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid master password",
		})
		return
	}

	err = vault.UpdatePassword(req.Name, req.Password, req.PasswordConfirm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.VaultService.SaveVaultToRedis(ctx, vault)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

type deletePasswordRequest struct {
	VaultName      string `form:"vault_name" binding:"required"`
	MasterPassword string `form:"master_password" binding:"required"`
	PasswordName   string `form:"password_name" binding:"required"`
}

func (server *Server) deletePassword(ctx *gin.Context) {
	var req deletePasswordRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	vault, err := server.VaultService.LoadVaultFromRedis(ctx, req.VaultName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Vault not found",
		})
		return
	}

	if !server.VaultService.VerifyMasterPassword(ctx, req.VaultName, req.MasterPassword) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid master password",
		})
		return
	}

	err = vault.DeletePassword(req.PasswordName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.VaultService.SaveVaultToRedis(ctx, vault)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Password deleted successfully",
	})
}
