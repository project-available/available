package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/project-available/available.git/db/sqlc"
)

// server http request
type Server struct {
	store db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server{
	server := &Server{store: store}
	router := gin.Default()
    //account
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/accounts/update/:id", server.updateAccount)
	router.DELETE("/accounts/delete/:id", server.deleteAccount)

	//booking
	router.POST("/bookings", server.createBooking)
	router.GET("/bookings/:account_id", server.getBookingOfAccount)
	router.GET("/bookings", server.listBookings)
	router.POST("/bookings/update/:id", server.updateBooking)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorMessage(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}