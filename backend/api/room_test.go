package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/project-available/available/db/mock"
	db "github.com/project-available/available/db/sqlc"
	"github.com/project-available/available/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateRoom(t *testing.T) {
	room := randomRoom()

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"name":     room.Name,
				"location": room.Location,
				"image":    room.Image,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateRoom(gomock.Any(), gomock.Any()).
					Times(1).
					Return(room, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidJSON",
			body: map[string]interface{}{
				"name": room.Name,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateRoom(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: map[string]interface{}{
				"name":     room.Name,
				"location": room.Location,
				"image":    room.Image,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateRoom(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Room{}, sql.ErrConnDone)
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

			url := "/rooms"
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

func TestGetRoom(t *testing.T) {
	room := randomRoom()
	customFields := []db.CustomField{
		{
			ID:    1,
			Key:   "capacity",
			Shown: true,
		},
	}
	customFieldValues := []db.CustomFieldsValue{
		{
			ID:            1,
			RoomID:        room.ID,
			CustomfieldID: 1,
			Value:         "20",
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
			roomID: room.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoom(gomock.Any(), gomock.Eq(room.ID)).
					Times(1).
					Return(room, nil)
				store.EXPECT().
					ListCustomFields(gomock.Any()).
					Times(1).
					Return(customFields, nil)
				store.EXPECT().
					ListRoomCustomFieldValues(gomock.Any(), gomock.Eq(room.ID)).
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
					GetRoom(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "RoomNotFound",
			roomID: room.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoom(gomock.Any(), gomock.Eq(room.ID)).
					Times(1).
					Return(db.Room{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "CustomFieldsError",
			roomID: room.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoom(gomock.Any(), gomock.Eq(room.ID)).
					Times(1).
					Return(room, nil)
				store.EXPECT().
					ListCustomFields(gomock.Any()).
					Times(1).
					Return([]db.CustomField{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "CustomFieldValuesError",
			roomID: room.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetRoom(gomock.Any(), gomock.Eq(room.ID)).
					Times(1).
					Return(room, nil)
				store.EXPECT().
					ListCustomFields(gomock.Any()).
					Times(1).
					Return(customFields, nil)
				store.EXPECT().
					ListRoomCustomFieldValues(gomock.Any(), gomock.Eq(room.ID)).
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

			url := fmt.Sprintf("/rooms/%d", tc.roomID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListRooms(t *testing.T) {
	n := 5
	rooms := make([]db.ListRoomsWithStatusRow, n)
	for i := 0; i < n; i++ {
		rooms[i] = db.ListRoomsWithStatusRow{
			ID:            utils.RandomInt(1, 1000),
			Location:      utils.RandomString(12),
			Name:          utils.RandomStringNumber(10),
			Image:         utils.RandomString(20),
			IsOccupiedNow: false,
			OccupiedUntil: nil,
		}
	}

	customFieldValues := []db.CustomFieldsValue{
		{
			ID:            1,
			RoomID:        rooms[0].ID,
			CustomfieldID: 1,
			Value:         "20",
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
					ListRoomsWithStatus(gomock.Any(), gomock.Any()).
					Times(1).
					Return(rooms, nil)
				store.EXPECT().
					ListRoomCustomFieldValuesBatch(gomock.Any(), gomock.Any()).
					Times(1).
					Return(customFieldValues, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "OKWithOccupiedRoom",
			query: map[string]string{
				"page_id":   "1",
				"page_size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				occupiedRooms := make([]db.ListRoomsWithStatusRow, 1)
				occupiedRooms[0] = db.ListRoomsWithStatusRow{
					ID:            utils.RandomInt(1, 1000),
					Location:      utils.RandomString(12),
					Name:          utils.RandomStringNumber(10),
					Image:         utils.RandomString(20),
					IsOccupiedNow: true,
					OccupiedUntil: time.Now().Add(2 * time.Hour),
				}
				store.EXPECT().
					ListRoomsWithStatus(gomock.Any(), gomock.Any()).
					Times(1).
					Return(occupiedRooms, nil)
				store.EXPECT().
					ListRoomCustomFieldValuesBatch(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.CustomFieldsValue{}, nil)
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
					ListRoomsWithStatus(gomock.Any(), gomock.Any()).
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
				"page_size": "200",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoomsWithStatus(gomock.Any(), gomock.Any()).
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
				"page_size": "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListRoomsWithStatus(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.ListRoomsWithStatusRow{}, sql.ErrConnDone)
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

			url := "/rooms"
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

func TestUpdateRoom(t *testing.T) {
	room := randomRoom()

	testCases := []struct {
		name          string
		roomID        int64
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			roomID: room.ID,
			body: map[string]interface{}{
				"name":     utils.RandomString(10),
				"location": utils.RandomString(10),
				"image":    utils.RandomString(20),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoom(gomock.Any(), gomock.Any()).
					Times(1).
					Return(room, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			roomID: 0,
			body: map[string]interface{}{
				"name":     utils.RandomString(10),
				"location": utils.RandomString(10),
				"image":    utils.RandomString(20),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoom(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InvalidBody",
			roomID: room.ID,
			body: map[string]interface{}{
				"name": utils.RandomString(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoom(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			roomID: room.ID,
			body: map[string]interface{}{
				"name":     utils.RandomString(10),
				"location": utils.RandomString(10),
				"image":    utils.RandomString(20),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoom(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Room{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			roomID: room.ID,
			body: map[string]interface{}{
				"name":     utils.RandomString(10),
				"location": utils.RandomString(10),
				"image":    utils.RandomString(20),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateRoom(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Room{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/rooms/update/%d", tc.roomID)
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

func TestDeleteRoom(t *testing.T) {
	room := randomRoom()

	testCases := []struct {
		name          string
		roomID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			roomID: room.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRoom(gomock.Any(), gomock.Eq(room.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			roomID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRoom(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			roomID: room.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRoom(gomock.Any(), gomock.Eq(room.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			roomID: room.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteRoom(gomock.Any(), gomock.Eq(room.ID)).
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

			url := fmt.Sprintf("/rooms/delete/%d", tc.roomID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			token, err := server.tokenMaker.CreateToken("testuser", server.config.AccessTokenDuration)
			require.NoError(t, err)
			request.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}