package db

import (
	"context"
	"testing"
	"time"

	"github.com/project-available/available.git/utils"
	"github.com/stretchr/testify/require"
)



func createRandomBooking(t *testing.T) Booking {
	account := createRandomAccount(t)
	room := createRandomRoom(t)
	phoneBooking := utils.RandomPhone()
	start := time.Now().UTC()
    end := start.Add(time.Hour)

	arg := CreateBookingParams{
		AccountID: account.ID,
		RoomID:    room.ID,
		Start: start,
		End:   end,
		PhoneBooking: phoneBooking,
	}
	booking, err := testQuery.CreateBooking(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, booking)
	require.Equal(t, arg.AccountID, booking.AccountID)
	require.Equal(t, arg.RoomID, booking.RoomID)
	require.WithinDuration(t, arg.Start, booking.Start, time.Millisecond)
	require.WithinDuration(t, arg.End, booking.End, time.Millisecond)
	require.NotZero(t, booking.ID)
	return booking
}
func TestCreateBooking(t *testing.T) {
    createRandomBooking(t)
}

func TestGetBookingOfAccount(t *testing.T) {
	booking1 := createRandomBooking(t)
	bookings, err := testQuery.GetBookingOfAccount(context.Background(), booking1.AccountID)
    require.NoError(t, err)
    require.NotEmpty(t, bookings)

    booking2 := bookings[0]
	require.Equal(t, booking1.ID, booking2.ID)
	require.Equal(t, booking1.AccountID, booking2.AccountID)
	require.Equal(t, booking1.RoomID, booking2.RoomID)
	require.WithinDuration(t, booking1.Start, booking2.Start, time.Millisecond)
	require.WithinDuration(t, booking1.End, booking2.End, time.Millisecond)
}

func TestListBookings(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomBooking(t)
	}	
	arg := ListBookingsParams{
		Limit:  5,
		Offset: 5,
	}
	bookings, err := testQuery.ListBookings(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, bookings, 5)
	for _, booking := range bookings {
		require.NotEmpty(t, booking)
	}	
}

func TestUpdateBooking(t *testing.T) {
	booking1 := createRandomBooking(t)

	newStatus := utils.RandomStatus()

	arg := UpdateBookingParams{
		ID:     booking1.ID,
		Status: newStatus,
	}
	booking2, err := testQuery.UpdateBooking(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, booking2)

	require.Equal(t, newStatus, booking2.Status)

	require.Equal(t, booking1.ID, booking2.ID)
	require.Equal(t, booking1.AccountID, booking2.AccountID)
	require.Equal(t, booking1.RoomID, booking2.RoomID)
	require.WithinDuration(t, booking1.Start, booking2.Start, time.Millisecond)
	require.WithinDuration(t, booking1.End, booking2.End, time.Millisecond)
}