package api

import (
	db "github.com/bytepharoh/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Create a server serves requests of our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// Newserver that creates a new http server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.POST("/transfers", server.createTransfer)
	router.GET("/accounts/:id", server.getAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.GET("/accounts", server.listAccount)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
