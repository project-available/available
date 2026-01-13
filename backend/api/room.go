package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/project-available/available/db/sqlc"
	"github.com/project-available/available/utils"
)

type CreateRoomRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	Image    string `json:"image" binding:"required"`
}

// CreateRoom godoc
// @Summary Create a new room
// @Description Create a new room (admin only)
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body CreateRoomRequest true "Room details"
// @Success 200 {object} db.Room
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /rooms [post]
func (server *Server) createRoom(ctx *gin.Context) {
	var req CreateRoomRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	arg := db.CreateRoomParams{
		Name:     req.Name,
		Location: req.Location,
		Image:    req.Image,
	}
	room, err := server.store.CreateRoom(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, room)
}

type GetRoomRequest struct {
	RoomID int64 `uri:"room_id" binding:"required"`
}
type CustomFieldResp struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Shown bool   `json:"shown"`
}
type GetRoomResponse struct {
	RoomID       int64             `json:"roomId"`
	Location     string            `json:"location"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	CustomFields []CustomFieldResp `json:"customFields"`
}

// GetRoom godoc
// @Summary Get room details
// @Description Get detailed information about a specific room including custom fields
// @Tags rooms
// @Accept json
// @Produce json
// @Param room_id path int true "Room ID"
// @Success 200 {object} GetRoomResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rooms/{room_id} [get]
func (server *Server) getRoom(ctx *gin.Context) {
	var req GetRoomRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	room, err := server.store.GetRoom(ctx, req.RoomID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	// Fetch all customfields
	defs, err := server.store.ListCustomFields(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	// Fetch all values for this room
	vals, err := server.store.ListRoomCustomFieldValues(ctx, req.RoomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	valMap := make(map[int64]string)
	for _, v := range vals {
		valMap[v.CustomfieldID] = v.Value
	}

	var customfields []CustomFieldResp
	for _, def := range defs {
		value := ""
		if v, ok := valMap[def.ID]; ok {
			value = v
		}
		customfields = append(customfields, CustomFieldResp{
			Key:   def.Key,
			Value: value,
			Shown: def.Shown,
		})
	}

	resp := GetRoomResponse{
		RoomID:       room.ID,
		Location:     room.Location,
		Name:         room.Name,
		Image:        room.Image,
		CustomFields: customfields,
	}
	ctx.JSON(http.StatusOK, resp)
}

type ListRoomsRequest struct {
	PageID   int32 `form:"page_id,required" binding:"min=1"`
	PageSize int32 `form:"page_size,required" binding:"min=1,max=100"`
}

type ListRoomResponse struct {
	RoomID       int64     `json:"roomId"`
	Status       string    `json:"status"`
	AvailableAt  time.Time `json:"availableAt"`
	Location     string    `json:"location"`
	Name         string    `json:"name"`
	Image        string    `json:"image"`
	CustomFields []string  `json:"customFields"`
}

// ListRooms godoc
// @Summary List all rooms
// @Description Get a paginated list of all rooms with their current status
// @Tags rooms
// @Accept json
// @Produce json
// @Param page_id query int true "Page number" minimum(1)
// @Param page_size query int true "Page size" minimum(1) maximum(100)
// @Success 200 {array} ListRoomResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rooms [get]
func (server *Server) listRooms(ctx *gin.Context) {
	var req ListRoomsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	t := time.Now()
	arg := db.ListRoomsWithStatusParams{
		Start:  t,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	rooms, err := server.store.ListRoomsWithStatus(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	// Extract room IDs for batch fetching custom fields
	roomIDs := make([]int64, 0, len(rooms))
	for _, room := range rooms {
		roomIDs = append(roomIDs, room.ID)
	}

	// Batch fetch all custom field values for all rooms in one query
	customFields, err := server.store.ListRoomCustomFieldValuesBatch(ctx, roomIDs)
	if err != nil {
		log.Default().Println(err)
	}

	// Group custom field values by room_id
	customFieldsByRoom := make(map[int64][]string)
	for _, v := range customFields {
		customFieldsByRoom[v.RoomID] = append(customFieldsByRoom[v.RoomID], v.Value)
	}

	// Build response
	resq := make([]ListRoomResponse, 0, len(rooms))
	for _, room := range rooms {
		var availableAt time.Time
		status := utils.Available
		if room.IsOccupiedNow.(bool) {
			status = utils.Booked
			availableAt = room.OccupiedUntil.(time.Time)
		}

		item := ListRoomResponse{
			RoomID:       room.ID,
			Status:       string(status),
			AvailableAt:  availableAt,
			Location:     room.Location,
			Name:         room.Name,
			Image:        room.Image,
			CustomFields: customFieldsByRoom[room.ID],
		}
		if item.CustomFields == nil {
			item.CustomFields = []string{}
		}
		resq = append(resq, item)
	}
	ctx.JSON(http.StatusOK, resq)
}

type UpdateRoomUriRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type UpdateRoomJsonRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	Image    string `json:"image" binding:"required"`
}

// UpdateRoom godoc
// @Summary Update room
// @Description Update room information (admin only)
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Param room body UpdateRoomJsonRequest true "Room update details"
// @Success 200 {object} db.Room
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /rooms/update/{id} [post]
func (server *Server) updateRoom(ctx *gin.Context) {
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
		ID:       uriReq.ID,
		Name:     jsonReq.Name,
		Location: jsonReq.Location,
		Image:    jsonReq.Image,
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

// DeleteRoom godoc
// @Summary Delete room
// @Description Soft delete a room (admin only)
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /rooms/delete/{id} [delete]
func (server *Server) deleteRoom(ctx *gin.Context) {
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
