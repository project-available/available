package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/project-available/available.git/db/sqlc"
)

// server http request
type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}

	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()
	//account
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:student_id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:student_id", server.deleteAccount)

	//booking
	router.POST("/bookings", server.createBooking)
	router.GET("/bookings/:account_id", server.getBookingOfAccount)
	router.GET("/bookings", server.listBookings)
	router.POST("/bookings/update/:id", server.updateBooking)

	//room
	router.POST("/rooms", server.createRoom)
	router.GET("/rooms", server.listRooms)
	router.POST("/rooms/update/:id", server.updateRoom)
	router.DELETE("/rooms/delete/:id", server.deleteRoom)
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
