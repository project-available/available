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
	mockdb "github.com/project-available/available/db/mock"
	db "github.com/project-available/available/db/sqlc"
	"github.com/project-available/available/utils"
	"github.com/stretchr/testify/require"
)

var now = time.Now()

func TestCreateBooking(t *testing.T) {
	account := randomUserAccount()
	room := randomRoom()
	booking := db.Booking{
		ID:           utils.RandomInt(1, 1000),
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        now,
		End:          now.Add(2 * time.Hour),
		Status:       "confirmed",
		PhoneBooking: account.Phone,
	}

	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)

	expectedResult := db.BookingTxResult{
		Booking: booking,
	}
	store.EXPECT().BookingTx(gomock.Any(), gomock.Any()).Times(1).Return(expectedResult, nil)

	// start test
	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := "/bookings"
	body := CreateBookingRequest{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        booking.Start,
		End:          booking.End,
		PhoneBooking: account.Phone,
	}
	data, err := json.Marshal(body)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	token, err := server.tokenMaker.CreateToken(account.StudentID, server.config.AccessTokenDuration)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+token)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomUserAccount() db.Account {
	return db.Account{
		ID:             utils.RandomInt(1, 1000),
		Name:           utils.RandomString(10),
		Role:           "user",
		Email:          utils.RandomEmail(),
		HashedPassword: utils.RandomString(20),
		Phone:          utils.RandomPhone(),
		StudentID:      utils.RandomStringNumber(7),
		IsDelete:       false,
	}
}

func randomRoom() db.Room {
	return db.Room{
		ID:       utils.RandomInt(1, 1000),
		Location: utils.RandomString(12),
		Name:     utils.RandomStringNumber(10),
		Image:    utils.RandomString(20),
		IsDelete: false,
	}
}

func TestCreateBookingOverlapConflict(t *testing.T) {
	account := randomUserAccount()
	room := randomRoom()

	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)

	startTime := now
	endTime := startTime.Add(2 * time.Hour)

	store.EXPECT().BookingTx(gomock.Any(), gomock.Any()).Times(1).
		Return(db.BookingTxResult{}, errors.New("booking time overlaps with existing bookings on this room"))

	// start test
	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := "/bookings"
	body := CreateBookingRequest{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        startTime,
		End:          endTime,
		PhoneBooking: account.Phone,
	}
	data, err := json.Marshal(body)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	token, err := server.tokenMaker.CreateToken(account.StudentID, server.config.AccessTokenDuration)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+token)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestCreateBookingInvalidTime(t *testing.T) {
	account := randomUserAccount()
	room := randomRoom()

	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().BookingTx(gomock.Any(), gomock.Any()).Times(0)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := "/bookings"
	body := CreateBookingRequest{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        now.Add(2 * time.Hour),
		End:          now,
		PhoneBooking: account.Phone,
	}
	data, err := json.Marshal(body)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	token, err := server.tokenMaker.CreateToken(account.StudentID, server.config.AccessTokenDuration)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+token)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCreateBookingInvalidJSON(t *testing.T) {
	account := randomUserAccount()

	ctrl := gomock.NewController(t)
	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().BookingTx(gomock.Any(), gomock.Any()).Times(0)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := "/bookings"
	body := map[string]interface{}{
		"account_id": account.ID,
	}
	data, err := json.Marshal(body)
	require.NoError(t, err)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	token, err := server.tokenMaker.CreateToken(account.StudentID, server.config.AccessTokenDuration)
	require.NoError(t, err)
	request.Header.Set("Authorization", "Bearer "+token)

	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetBookingOfAccount(t *testing.T) {
	account := randomUserAccount()
	bookings := []db.Booking{
		{
			ID:           utils.RandomInt(1, 1000),
			AccountID:    account.ID,
			RoomID:       utils.RandomInt(1, 100),
			Start:        now,
			End:          now.Add(2 * time.Hour),
			Status:       "confirmed",
			PhoneBooking: account.Phone,
		},
	}

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBookingOfAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(bookings, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "InvalidAccountID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBookingOfAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBookingOfAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return([]db.Booking{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBookingOfAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return([]db.Booking{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/bookings/%d", tc.accountID)
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

func TestListBookings(t *testing.T) {
	bookings := []db.Booking{
		{
			ID:           utils.RandomInt(1, 1000),
			AccountID:    utils.RandomInt(1, 100),
			RoomID:       utils.RandomInt(1, 100),
			Start:        now,
			End:          now.Add(2 * time.Hour),
			Status:       "confirmed",
			PhoneBooking: utils.RandomPhone(),
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
					ListBookings(gomock.Any(), gomock.Any()).
					Times(1).
					Return(bookings, nil)
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
					ListBookings(gomock.Any(), gomock.Any()).
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
				"page_size": "10",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListBookings(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Booking{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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
					ListBookings(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Booking{}, sql.ErrConnDone)
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

			url := "/bookings"
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

func TestUpdateBooking(t *testing.T) {
	account := randomUserAccount()
	booking := db.Booking{
		ID:           utils.RandomInt(1, 1000),
		AccountID:    account.ID,
		RoomID:       utils.RandomInt(1, 100),
		Start:        now,
		End:          now.Add(2 * time.Hour),
		Status:       "confirmed",
		PhoneBooking: account.Phone,
	}

	testCases := []struct {
		name          string
		bookingID     int64
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			bookingID: booking.ID,
			body: map[string]interface{}{
				"status": "confirmed",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateBooking(gomock.Any(), gomock.Any()).
					Times(1).
					Return(booking, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "InvalidBookingID",
			bookingID: 0,
			body: map[string]interface{}{
				"status": "confirmed",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateBooking(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InvalidStatus",
			bookingID: booking.ID,
			body: map[string]interface{}{
				"status": "invalid_status",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateBooking(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InvalidJSON",
			bookingID: booking.ID,
			body:      map[string]interface{}{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateBooking(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			bookingID: booking.ID,
			body: map[string]interface{}{
				"status": "confirmed",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateBooking(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Booking{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			bookingID: booking.ID,
			body: map[string]interface{}{
				"status": "confirmed",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateBooking(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Booking{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/bookings/update/%d", tc.bookingID)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			token, err := server.tokenMaker.CreateToken(account.StudentID, server.config.AccessTokenDuration)
			require.NoError(t, err)
			request.Header.Set("Authorization", "Bearer "+token)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetRoomBookings(t *testing.T) {
	room := randomRoom()
	bookings := []db.Booking{
		{
			ID:           utils.RandomInt(1, 1000),
			AccountID:    utils.RandomInt(1, 100),
			RoomID:       room.ID,
			Start:        now,
			End:          now.Add(2 * time.Hour),
			Status:       "confirmed",
			PhoneBooking: utils.RandomPhone(),
		},
	}

	testCases := []struct {
		name          string
		roomID        int64
		date          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			roomID: room.ID,
			date:   "2024-01-01",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBookingsOnDate(gomock.Any(), gomock.Any()).
					Times(1).
					Return(bookings, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "InvalidRoomID",
			roomID: 0,
			date:   "2024-01-01",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBookingsOnDate(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InvalidDate",
			roomID: room.ID,
			date:   "invalid-date",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBookingsOnDate(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "MissingDate",
			roomID: room.ID,
			date:   "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBookingsOnDate(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			roomID: room.ID,
			date:   "2024-01-01",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetBookingsOnDate(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Booking{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/rooms/%d/bookings", tc.roomID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			if tc.date != "" {
				q := request.URL.Query()
				q.Add("date", tc.date)
				request.URL.RawQuery = q.Encode()
			}

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
