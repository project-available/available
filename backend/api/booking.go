package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/project-available/available.git/db/sqlc"
)

type CreateBookingRequest struct {
	AccountID    int64     `json:"account_id" binding:"required"`
	RoomID       int64     `json:"room_id" binding:"required"`
	Start        time.Time `json:"start" binding:"required"`
	End          time.Time `json:"end" binding:"required"`
	PhoneBooking string    `json:"phone_booking" binding:"required"`
}

func (server *Server) createBooking(ctx *gin.Context) {
	var req CreateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	arg := db.CreateBookingParams{
		AccountID:    req.AccountID,
		RoomID:       req.RoomID,
		Start:        req.Start,
		End:          req.End,
		PhoneBooking: req.PhoneBooking,
	}
	booking, err := server.store.CreateBooking(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, booking)
}

type GetBookingOfAccountRequest struct {
	AccountID int64 `uri:"account_id" binding:"required"`
}

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
	PageID  int32 `form:"page_id,default=1" binding:"min=1"`
	PageSize int32 `form:"page_size,default=10" binding:"min=1,max=100"`
}

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