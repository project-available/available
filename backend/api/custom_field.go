package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/project-available/available/db/sqlc"
)

type CreateCustomFieldRequest struct {
	Key string `json:"key" binding:"required"`
}

// CreateCustomField godoc
// @Summary Create a new custom field
// @Description Create a new custom field definition (admin only)
// @Tags custom_fields
// @Accept json
// @Produce json
// @Param custom_field body CreateCustomFieldRequest true "Custom field details"
// @Success 200 {object} db.CustomField
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /custom_fields [post]
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

// ListCustomFields godoc
// @Summary List all custom fields
// @Description Get a list of all custom field definitions
// @Tags custom_fields
// @Accept json
// @Produce json
// @Success 200 {array} db.CustomField
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /custom_fields [get]
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

// ListRoomCustomFieldValues godoc
// @Summary Get custom field values for a room
// @Description Get all custom field values for a specific room
// @Tags custom_fields
// @Accept json
// @Produce json
// @Param room_id path int true "Room ID"
// @Success 200 {array} db.CustomFieldsValue
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /custom_fields/{room_id} [get]
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

// UpsertRoomCustomFieldValue godoc
// @Summary Update custom field value
// @Description Create or update a custom field value for a room (admin only)
// @Tags custom_fields
// @Accept json
// @Produce json
// @Param value body UpsertRoomCustomFieldValueRequest true "Custom field value"
// @Success 200 {string} string "success"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /custom_fields [put]
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
