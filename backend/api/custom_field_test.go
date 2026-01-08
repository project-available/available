package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/project-available/available/db/mock"
	db "github.com/project-available/available/db/sqlc"
	"github.com/project-available/available/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateCustomField(t *testing.T) {
	customField := db.CustomField{
		ID:    utils.RandomInt(1, 1000),
		Key:   utils.RandomString(10),
		Shown: true,
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
				"key": customField.Key,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCustomField(gomock.Any(), gomock.Eq(customField.Key)).
					Times(1).
					Return(customField, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidJSON",
			body: map[string]interface{}{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCustomField(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: map[string]interface{}{
				"key": customField.Key,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCustomField(gomock.Any(), gomock.Eq(customField.Key)).
					Times(1).
					Return(db.CustomField{}, sql.ErrConnDone)
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

			url := "/custom_fields"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			token, err := server.tokenMaker.CreateToken("testuser", server.config.AccessTokenDuration)
			require.NoError(t, err)
			request.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListCustomFields(t *testing.T) {
	customFields := []db.CustomField{
		{
			ID:    utils.RandomInt(1, 1000),
			Key:   utils.RandomString(10),
			Shown: true,
		},
		{
			ID:    utils.RandomInt(1, 1000),
			Key:   utils.RandomString(10),
			Shown: false,
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
				"page_size": "10",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCustomFields(gomock.Any()).
					Times(1).
					Return(customFields, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: map[string]string{
				"page_id":   "0",
				"page_size": "10",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCustomFields(gomock.Any()).
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
					ListCustomFields(gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: map[string]string{
				"page_id":   "1",
				"page_size": "10",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCustomFields(gomock.Any()).
					Times(1).
					Return([]db.CustomField{}, sql.ErrConnDone)
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

			url := "/custom_fields"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			for key, value := range tc.query {
				q.Add(key, value)
			}
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListRoomCustomFieldValues(t *testing.T) {
	roomID := utils.RandomInt(1, 1000)
	customFieldValues := []db.CustomFieldsValue{
		{
			ID:            utils.RandomInt(1, 1000),
			RoomID:        roomID,
			CustomfieldID: utils.RandomInt(1, 100),
			Value:         utils.RandomString(10),
		},
		{
			ID:            utils.RandomInt(1, 1000),
			RoomID:        roomID,
			CustomfieldID: utils.RandomInt(1, 100),
			Value:         utils.RandomString(10),
		},
	}

	testCases := []struct {
		name          string
		roomID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			roomID: roomID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoomCustomFieldValues(gomock.Any(), gomock.Eq(roomID)).
					Times(1).
					Return(customFieldValues, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InvalidRoomID",
			roomID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoomCustomFieldValues(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			roomID: roomID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoomCustomFieldValues(gomock.Any(), gomock.Eq(roomID)).
					Times(1).
					Return([]db.CustomFieldsValue{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/custom_fields/%d", tc.roomID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpsertRoomCustomFieldValue(t *testing.T) {
	roomID := utils.RandomInt(1, 1000)
	customFieldID := utils.RandomInt(1, 100)
	value := utils.RandomString(10)

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"room_id":         roomID,
				"custom_field_id": customFieldID,
				"value":           value,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpsertRoomCustomFieldValue(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidJSON",
			body: map[string]interface{}{
				"room_id": roomID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpsertRoomCustomFieldValue(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidRoomID",
			body: map[string]interface{}{
				"room_id":         0,
				"custom_field_id": customFieldID,
				"value":           value,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpsertRoomCustomFieldValue(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidCustomFieldID",
			body: map[string]interface{}{
				"room_id":         roomID,
				"custom_field_id": 0,
				"value":           value,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpsertRoomCustomFieldValue(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: map[string]interface{}{
				"room_id":         roomID,
				"custom_field_id": customFieldID,
				"value":           value,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpsertRoomCustomFieldValue(gomock.Any(), gomock.Any()).
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/custom_fields"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			token, err := server.tokenMaker.CreateToken("testuser", server.config.AccessTokenDuration)
			require.NoError(t, err)
			request.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}