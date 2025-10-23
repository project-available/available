package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/project-available/available.git/db/sqlc"
	"github.com/project-available/available.git/token"
	"github.com/project-available/available.git/utils"
)

// server http request
type Server struct {
	config     utils.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w, %d", err, len(config.TokenSymmetricKey))
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.POST("/accounts/login", server.loginAccount)

	router.GET("/rooms", server.listRooms)

	router.GET("/bookings", server.listBookings)

	authRoutes := router.Group("/").Use(authMiddleWare(server.tokenMaker))

	authRoutes.GET("/accounts/:student_id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PUT("/accounts/:id", server.updateAccount)
	authRoutes.DELETE("/accounts/:student_id", server.deleteAccount)

	authRoutes.GET("/bookings/:account_id", server.getBookingOfAccount)
	authRoutes.POST("/bookings/update/:id", server.updateBooking)
	authRoutes.POST("/bookings", server.createBooking)

	authRoutes.POST("/rooms", server.createRoom)
	authRoutes.POST("/rooms/update/:id", server.updateRoom)
	authRoutes.DELETE("/rooms/delete/:id", server.deleteRoom)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorMessage(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
