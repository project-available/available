package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	mockdb "github.com/project-available/available/db/mock"
	db "github.com/project-available/available/db/sqlc"
	"github.com/project-available/available/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	password := utils.RandomString(10)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	account := db.Account{
		ID:             utils.RandomInt(1, 1000),
		Name:           utils.RandomString(10),
		Role:           string(utils.RoleUser),
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
		Phone:          utils.RandomPhone(),
		StudentID:      utils.RandomStudentID(),
		IsDelete:       false,
	}

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"name":       account.Name,
				"email":      account.Email,
				"password":   password,
				"phone":      account.Phone,
				"student_id": account.StudentID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidJSON",
			body: map[string]interface{}{
				"name": account.Name,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UniqueViolation",
			body: map[string]interface{}{
				"name":       account.Name,
				"email":      account.Email,
				"password":   password,
				"phone":      account.Phone,
				"student_id": account.StudentID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: map[string]interface{}{
				"name":       account.Name,
				"email":      account.Email,
				"password":   password,
				"phone":      account.Phone,
				"student_id": account.StudentID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetAccount(t *testing.T) {
	account := randomUserAccount()

	testCases := []struct {
		name          string
		studentID     string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			studentID: account.StudentID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.StudentID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			studentID: account.StudentID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.StudentID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			studentID: account.StudentID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.StudentID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%s", tc.studentID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			token, err := server.tokenMaker.CreateToken(account.StudentID, server.config.AccessTokenDuration)
			require.NoError(t, err)
			request.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListAccounts(t *testing.T) {
	accounts := []db.ListAccountsRow{
		{
			ID:        utils.RandomInt(1, 1000),
			Name:      utils.RandomString(10),
			Role:      string(utils.RoleUser),
			Email:     utils.RandomEmail(),
			Phone:     utils.RandomPhone(),
			StudentID: utils.RandomStudentID(),
		},
		{
			ID:        utils.RandomInt(1, 1000),
			Name:      utils.RandomString(10),
			Role:      string(utils.RoleUser),
			Email:     utils.RandomEmail(),
			Phone:     utils.RandomPhone(),
			StudentID: utils.RandomStudentID(),
		},
	}

	testCases := []struct {
		name          string
		query         map[string]string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: map[string]string{
				"page_id":   "1",
				"page_size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: map[string]string{
				"page_id":   "0",
				"page_size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: map[string]string{
				"page_id":   "1",
				"page_size": "3",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			query: map[string]string{
				"page_id":   "1",
				"page_size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListAccountsRow{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: map[string]string{
				"page_id":   "1",
				"page_size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListAccountsRow{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			for key, value := range tc.query {
				q.Add(key, value)
			}
			request.URL.RawQuery = q.Encode()

			token, err := server.tokenMaker.CreateToken("testuser", server.config.AccessTokenDuration)
			require.NoError(t, err)
			request.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	account := randomUserAccount()

	testCases := []struct {
		name          string
		accountID     int64
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			body: map[string]interface{}{
				"name":  utils.RandomString(10),
				"phone": utils.RandomPhone(),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			body: map[string]interface{}{
				"name":  utils.RandomString(10),
				"phone": utils.RandomPhone(),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InvalidBody",
			accountID: account.ID,
			body: map[string]interface{}{
				"name": utils.RandomString(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			body: map[string]interface{}{
				"name":  utils.RandomString(10),
				"phone": utils.RandomPhone(),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			token, err := server.tokenMaker.CreateToken(account.StudentID, server.config.AccessTokenDuration)
			require.NoError(t, err)
			request.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	account := randomUserAccount()

	testCases := []struct {
		name          string
		studentID     string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			studentID: account.StudentID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.StudentID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			studentID: account.StudentID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.StudentID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			studentID: account.StudentID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.StudentID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%s", tc.studentID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			token, err := server.tokenMaker.CreateToken(account.StudentID, server.config.AccessTokenDuration)
			require.NoError(t, err)
			request.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestLoginAccount(t *testing.T) {
	password := utils.RandomString(10)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	account := db.Account{
		ID:             utils.RandomInt(1, 1000),
		Name:           utils.RandomString(10),
		Role:           string(utils.RoleUser),
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
		Phone:          utils.RandomPhone(),
		StudentID:      utils.RandomStudentID(),
		IsDelete:       false,
	}

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"email":    account.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountByEmail(gomock.Any(), gomock.Eq(account.Email)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidJSON",
			body: map[string]interface{}{
				"email": account.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountByEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			body: map[string]interface{}{
				"email":    account.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountByEmail(gomock.Any(), gomock.Eq(account.Email)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: map[string]interface{}{
				"email":    account.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountByEmail(gomock.Any(), gomock.Eq(account.Email)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "WrongPassword",
			body: map[string]interface{}{
				"email":    account.Email,
				"password": "wrongpassword",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountByEmail(gomock.Any(), gomock.Eq(account.Email)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	testCases := []struct {
		name      string
		config    utils.Config
		expectErr bool
	}{
		{
			name: "OK",
			config: utils.Config{
				TokenSymmetricKey:   utils.RandomString(32),
				AccessTokenDuration: time.Minute,
			},
			expectErr: false,
		},
		{
			name: "InvalidTokenKey",
			config: utils.Config{
				TokenSymmetricKey:   utils.RandomString(10),
				AccessTokenDuration: time.Minute,
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server, err := NewServer(tc.config, store)

			if tc.expectErr {
				require.Error(t, err)
				require.Nil(t, server)
			} else {
				require.NoError(t, err)
				require.NotNil(t, server)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	err := errors.New("test error")
	result := errorMessage(err)

	require.NotNil(t, result)
	require.Equal(t, "test error", result["error"])
}
