package api

import (
	"github.com/gin-contrib/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sordgom/PasswordManager/server/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config       config.Config
	router       *gin.Engine
	VaultService config.VaultService
}

func NewServer(conf config.Config, vaultService config.VaultService) (*Server, error) {

	server := &Server{
		config:       conf,
		VaultService: vaultService,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Handle("GET", "/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(cors.Default())

	// Vault API
	router.POST("/vault", server.createVault)

	// Password API
	router.POST("/password", server.createPassword)
	router.GET("/passwords", server.getPasswords)
	router.GET("/password", server.getPassword)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
