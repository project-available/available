package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/project-available/available.git/db/sqlc"
	"github.com/project-available/available.git/utils"
)

type createAccountRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	StudentID string `json:"student_id" binding:"required"`
}

type accountResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	StudentID string `json:"student_id"`
}

func newAccountResponse(user db.Account) accountResponse {
	return accountResponse{
		ID:        user.ID,
		Name:      user.Name,
		Role:      user.Role,
		Email:     user.Email,
		Phone:     user.Phone,
		StudentID: user.StudentID,
	}
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	arg := db.CreateAccountParams{
		Name:           req.Name,
		Role:           string(utils.RoleUser),
		Email:          req.Email,
		HashedPassword: hashedPassword,
		Phone:          req.Phone,
		StudentID:      req.StudentID,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorMessage(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	res := newAccountResponse(account)
	ctx.JSON(http.StatusOK, res)
}

type getAccountRequest struct {
	StudentID string `uri:"student_id" binding:"required"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	account, err := server.store.GetAccount(ctx, req.StudentID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	res := newAccountResponse(account)
	ctx.JSON(http.StatusOK, res)
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
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
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
		ID:    uriReq.ID,
		Name:  jsonReq.Name,
		Phone: jsonReq.Phone,
	}
	account, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type deleteAccountRequest struct {
	StudentID string `uri:"student_id" binding:"required"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}
	err := server.store.DeleteAccount(ctx, req.StudentID)
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

type loginAccountRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginAccountResponse struct {
	AccessToken string          `json:"access_token"`
	Account     accountResponse `json:"account"`
}

func (server *Server) loginAccount(ctx *gin.Context) {
	var req loginAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorMessage(err))
		return
	}

	account, err := server.store.GetAccountByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorMessage(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	err = utils.CheckPassword(req.Password, account.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorMessage(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(req.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorMessage(err))
		return
	}

	res := loginAccountResponse{
		AccessToken: accessToken,
		Account:     newAccountResponse(account),
	}

	ctx.JSON(http.StatusOK, res)
}
