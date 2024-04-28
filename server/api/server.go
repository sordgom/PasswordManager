package api

import (
	"github.com/sordgom/PasswordManager/server/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	config       config.Config
	router       *gin.Engine
	VaultService *config.RedisVaultService
}

func NewServer(conf config.Config) (*Server, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddress,
		Password: "", // no password set
		DB:       0,  // default db
	})

	vaultService := config.NewRedisVaultService(client)

	server := &Server{
		config:       conf,
		VaultService: vaultService,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	// Vault API
	router.POST("/vault", server.createVault)

	// Password API
	router.POST("/password", server.createPassword)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
