package api

import (
	"fmt"

	db "github.com/bytepharoh/simplebank/db/sqlc"
	"github.com/bytepharoh/simplebank/toeken"
	"github.com/bytepharoh/simplebank/util"
	"github.com/gin-gonic/gin"
)

// Create a server serves requests of our banking service
type Server struct {
	store      db.Store
	tokenMaker toeken.Maker
	config     util.Config
	router     *gin.Engine
}

// Newserver that creates a new http server and setup routing
func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := toeken.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/transfers", server.createTransfer)

	router.GET("/accounts/:id", server.getAccount)
	router.GET("/users/:username", server.getUser)
	router.GET("/accounts", server.listAccount)

	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
