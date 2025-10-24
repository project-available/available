package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/project-available/available.git/db/sqlc"
)

type CreateRoomRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	Image    string `json:"image" binding:"required"`
}

func (server *Server) createRoom(ctx *gin.Context) {
	var req CreateRoomRequest

	if err := ctx.ShouldBindJSON(&req) ; err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	arg := db.CreateRoomParams{
		Name: req.Name,
		Location: req.Location,
		Image: req.Image,
	}
	room, err := server.store.CreateRoom(ctx,arg);
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, room)
}

type ListRoomsRequest struct {
	PageID  int32 `form:"page_id,default=1" binding:"min=1"`
	PageSize int32 `form:"page_size,default=10" binding:"min=1,max=100"`
}

func (server *Server) listRooms (ctx *gin.Context){
	var req ListRoomsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	arg := db.ListRoomsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	rooms, err := server.store.ListRooms(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, rooms)
}

type UpdateRoomUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type UpdateRoomJsonRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	Image    string `json:"image" binding:"required"`
}
func (server *Server) updateRoom(ctx *gin.Context){
	var uriReq UpdateRoomUriRequest
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	var jsonReq UpdateRoomJsonRequest
	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	arg := db.UpdateRoomParams{
		ID: uriReq.ID,
		Name: jsonReq.Name,
		Location: jsonReq.Location,
		Image: jsonReq.Image,
	}
	room, err := server.store.UpdateRoom(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, room)
}

type DeleteRoomRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
func (server *Server) deleteRoom(ctx *gin.Context){
	var req DeleteRoomRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	err := server.store.DeleteRoom(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}