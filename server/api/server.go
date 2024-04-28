package api

import (
	"server/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config config.Config
	router *gin.Engine
}

func NewServer(config config.Config) (*Server, error) {
	server := &Server{
		config: config,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Vault API
	server.router.GET("/vault", server.createVault)

	// Password API
	server.router.GET("/password", server.createPassword)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
