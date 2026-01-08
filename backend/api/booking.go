package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/project-available/available/db/sqlc"
)

type CreateBookingRequest struct {
	AccountID    int64     `json:"account_id" binding:"required"`
	RoomID       int64     `json:"room_id" binding:"required"`
	Start        time.Time `json:"start" binding:"required"`
	End          time.Time `json:"end" binding:"required"`
	PhoneBooking string    `json:"phone_booking" binding:"required"`
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a new room booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param booking body CreateBookingRequest true "Booking details"
// @Success 200 {object} db.BookingTxResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /bookings [post]
func (server *Server) createBooking(ctx *gin.Context) {
	var req CreateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	// Validate times: start must be before end
	if !req.Start.Before(req.End) {
		ctx.JSON(http.StatusBadRequest, errorMessage(errors.New("start must be before end")))
		return
	}

	arg := db.BookingTxParams{
		AccountID:    req.AccountID,
		RoomID:       req.RoomID,
		Start:        req.Start,
		End:          req.End,
		PhoneBooking: req.PhoneBooking,
	}
	result, err := server.store.BookingTx(ctx, db.BookingTxParams(arg))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

type GetBookingOfAccountRequest struct {
	AccountID int64 `uri:"account_id" binding:"required"`
}

// GetBookingOfAccount godoc
// @Summary Get bookings for an account
// @Description Get all bookings for a specific account
// @Tags bookings
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {array} db.Booking
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /bookings/{account_id} [get]
func (server *Server) getBookingOfAccount(ctx *gin.Context) {
	var req GetBookingOfAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	bookings, err := server.store.GetBookingOfAccount(ctx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, bookings)
}

type ListBookingsRequest struct {
	PageID   int32 `form:"page_id,default=1" binding:"min=1"`
	PageSize int32 `form:"page_size,default=10" binding:"min=1,max=100"`
}

// ListBookings godoc
// @Summary List all bookings
// @Description Get a paginated list of all bookings
// @Tags bookings
// @Accept json
// @Produce json
// @Param page_id query int false "Page number" default(1) minimum(1)
// @Param page_size query int false "Page size" default(10) minimum(1) maximum(100)
// @Success 200 {array} db.Booking
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /bookings [get]
func (server *Server) listBookings(ctx *gin.Context) {
	var req ListBookingsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	arg := db.ListBookingsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	bookings, err := server.store.ListBookings(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, bookings)
}

type UpdateBookingUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateBookingJsonRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed cancelled"`
}

// UpdateBooking godoc
// @Summary Update booking status
// @Description Update the status of a booking
// @Tags bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param booking body UpdateBookingJsonRequest true "Booking update details"
// @Success 200 {object} db.Booking
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /bookings/update/{id} [post]
func (server *Server) updateBooking(ctx *gin.Context) {
	// Bind URI parameter
	var uriReq UpdateBookingUriRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	// Bind JSON body
	var jsonReq UpdateBookingJsonRequest
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	arg := db.UpdateBookingParams{
		ID:     uriReq.ID,
		Status: jsonReq.Status,
	}

	booking, err := server.store.UpdateBooking(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, booking)
}

type GetRoomBookingsUriRequest struct {
	RoomID int64 `uri:"room_id" binding:"required"`
}

type GetRoomBookingQueryRequest struct {
	Date string `form:"date" binding:"required"`
}

// GetRoomBookings godoc
// @Summary Get room bookings for a date
// @Description Get all bookings for a specific room on a specific date
// @Tags rooms
// @Accept json
// @Produce json
// @Param room_id path int true "Room ID"
// @Param date query string true "Date in YYYY-MM-DD format" example("2024-01-15")
// @Success 200 {array} db.Booking
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rooms/{room_id}/bookings [get]
func (server *Server) getRoomBookings(ctx *gin.Context) {
	var req1 GetRoomBookingsUriRequest
	if err := ctx.ShouldBindUri(&req1); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	var req2 GetRoomBookingQueryRequest
	if err := ctx.ShouldBindQuery(&req2); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	// Parse only the date part (always safe)
	date, err := time.Parse("2006-01-02", req2.Date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	const timeZone = "Asia/Bangkok"
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		loc = time.FixedZone("ICT", 7*60*60)
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, loc)

	arg := db.GetBookingsOnDateParams{
		RoomID:  req1.RoomID,
		Start:   startOfDay,
		Start_2: endOfDay,
	}

	bookings, err := server.store.GetBookingsOnDate(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, bookings)
}
