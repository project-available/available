package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/project-available/available.git/db/sqlc"
)

type CreateRoomRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	Image    string `json:"image" binding:"required"`
}

// Handle room creation
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

// Handle getting room details info
// along with its custom fields and
// booking status of the room by day (default is now)
func (server *Server) getRoom(ctx *gin.Context) {
	var req GetRoomRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	room, err := server.store.GetRoom(ctx, req.RoomID)
	if err != nil {
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
	RoomID       int64    `json:"roomId"`
	Location     string   `json:"location"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	CustomFields []string `json:"customFields"`
}

func (server *Server) listRooms(ctx *gin.Context) {
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

	// Extract room IDs for batch fetching
	roomIDs := make([]int64, len(rooms))
	for i, room := range rooms {
		roomIDs[i] = room.ID
	}

	// Batch fetch all custom field values for all rooms in one query
	allCustomFieldValues, err := server.store.ListRoomCustomFieldValuesBatch(ctx, roomIDs)
	if err != nil {
		log.Default().Println(err)
	}

	// Group custom field values by room_id
	customFieldsByRoom := make(map[int64][]string)
	for _, v := range allCustomFieldValues {
		customFieldsByRoom[v.RoomID] = append(customFieldsByRoom[v.RoomID], v.Value)
	}

	// Build response
	resq := make([]ListRoomResponse, 0, len(rooms))
	for _, room := range rooms {
		item := ListRoomResponse{
			RoomID:       room.ID,
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
