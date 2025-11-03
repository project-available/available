package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/project-available/available.git/db/sqlc"
)

type CreateCustomFieldRequest struct {
	Key string `json:"key" binding:"required"`
}

func (server *Server) createCustomField(ctx *gin.Context) {
	var req CreateCustomFieldRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	customField, err := server.store.CreateCustomField(ctx, req.Key)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, customField)
}

type ListCustomFieldsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) listCustomFields(ctx *gin.Context) {
	var req ListCustomFieldsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	customFields, err := server.store.ListCustomFields(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, customFields)
}

type ListRoomCustomFieldValuesRequest struct {
	RoomID int64 `uri:"room_id" binding:"required,min=1"`
}

func (server *Server) listRoomCustomFieldValues(ctx *gin.Context) {
	var req ListRoomCustomFieldValuesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	customFieldValues, err := server.store.ListRoomCustomFieldValues(ctx, req.RoomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, customFieldValues)
}

type UpsertRoomCustomFieldValueRequest struct {
	RoomID        int64  `json:"room_id" binding:"required,min=1"`
	CustomFieldID int64  `json:"custom_field_id" binding:"required,min=1"`
	Value         string `json:"value" binding:"required"`
}

func (server *Server) upsertRoomCustomFieldValue(ctx *gin.Context) {
	var req UpsertRoomCustomFieldValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	arg := db.UpsertRoomCustomFieldValueParams{
		RoomID:        req.RoomID,
		CustomfieldID: req.CustomFieldID,
		Value:         req.Value,
	}

	err := server.store.UpsertRoomCustomFieldValue(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	ctx.JSON(http.StatusOK, "success")
}
