package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/project-available/available.git/db/sqlc"
)

type createAccountRequest struct {
	Name      string `json:"name" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=admin user"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	StudentID string `json:"student_id" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	arg := db.CreateAccountParams{
		Name:      req.Name,
		Role:      req.Role,
		Email:     req.Email,
		Password:  req.Password,
		Phone:     req.Phone,
		StudentID: req.StudentID,
	}
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

type updateAccountUriRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

type updateAccountJsonRequest struct {
    Name      string `json:"name" binding:"required"`
    Role      string `json:"role" binding:"required,oneof=admin user"`
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required"`
    Phone     string `json:"phone" binding:"required"`
    StudentID string `json:"student_id" binding:"required"`
}
func (server *Server) updateAccount(ctx *gin.Context) {
	var uriReq updateAccountUriRequest
    if err := ctx.ShouldBindUri(&uriReq); err != nil {
        ctx.JSON(http.StatusBadRequest, errorMessage(err))
        return
    }
	 var jsonReq updateAccountJsonRequest
    if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
        ctx.JSON(http.StatusBadRequest, errorMessage(err))
        return
    }
 arg := db.UpdateAccountParams{
        ID:        uriReq.ID,     
        Name:      jsonReq.Name,   
        Role:      jsonReq.Role,
        Email:     jsonReq.Email,
        Password:  jsonReq.Password,
        Phone:     jsonReq.Phone,
        StudentID: jsonReq.StudentID,
    }
	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}	
	account, err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}	
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}