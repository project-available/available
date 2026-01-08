package db

import (
	"context"
	"testing"
	"time"

	"github.com/project-available/available/utils"
	"github.com/stretchr/testify/require"
)

func TestBookingTx(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	room := createRandomRoom(t)
	start := time.Now().UTC()
	end := start.Add(time.Hour)
	phone := utils.RandomPhone()

	// 1. Successful Booking
	arg := BookingTxParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start,
		End:          end,
		PhoneBooking: phone,
	}

	result, err := store.BookingTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.Booking)
	require.Equal(t, account.ID, result.Booking.AccountID)
	require.Equal(t, room.ID, result.Booking.RoomID)
	require.Equal(t, "confirmed", result.Booking.Status)

	// 2. Overlapping Booking (should fail)
	arg2 := BookingTxParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start, // Same start
		End:          end,   // Same end
		PhoneBooking: phone,
	}
	result2, err := store.BookingTx(context.Background(), arg2)
	require.Error(t, err)
	require.EqualError(t, err, "booking time overlaps with existing bookings on this room")
	require.Empty(t, result2)

	// 3. Non-overlapping booking (should success)
	arg3 := BookingTxParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        end.Add(time.Minute),
		End:          end.Add(time.Hour),
		PhoneBooking: phone,
	}
	result3, err := store.BookingTx(context.Background(), arg3)
	require.NoError(t, err)
	require.NotEmpty(t, result3)
}
