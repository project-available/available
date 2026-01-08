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

func TestBookingTx_DifferentRooms(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	room1 := createRandomRoom(t)
	room2 := createRandomRoom(t)

	start := time.Now().UTC()
	end := start.Add(time.Hour)
	phone := utils.RandomPhone()

	arg1 := BookingTxParams{
		AccountID:    account.ID,
		RoomID:       room1.ID,
		Start:        start,
		End:          end,
		PhoneBooking: phone,
	}
	result1, err := store.BookingTx(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, result1)

	arg2 := BookingTxParams{
		AccountID:    account.ID,
		RoomID:       room2.ID,
		Start:        start,
		End:          end,
		PhoneBooking: phone,
	}
	result2, err := store.BookingTx(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, result2)
	require.NotEqual(t, result1.Booking.ID, result2.Booking.ID)
}

func TestBookingTx_PartialOverlap(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	room := createRandomRoom(t)

	start1 := time.Now().UTC()
	end1 := start1.Add(2 * time.Hour)
	phone := utils.RandomPhone()

	arg1 := BookingTxParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start1,
		End:          end1,
		PhoneBooking: phone,
	}
	result1, err := store.BookingTx(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, result1)

	start2 := start1.Add(time.Hour)
	end2 := end1.Add(time.Hour)

	arg2 := BookingTxParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start2,
		End:          end2,
		PhoneBooking: phone,
	}
	result2, err := store.BookingTx(context.Background(), arg2)
	require.Error(t, err)
	require.EqualError(t, err, "booking time overlaps with existing bookings on this room")
	require.Empty(t, result2.Booking.ID)
}

func TestBookingTx_AdjacentBookings(t *testing.T) {
	store := NewStore(testDB)

	account := createRandomAccount(t)
	room := createRandomRoom(t)

	start1 := time.Now().UTC()
	end1 := start1.Add(time.Hour)
	phone := utils.RandomPhone()

	arg1 := BookingTxParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start1,
		End:          end1,
		PhoneBooking: phone,
	}
	result1, err := store.BookingTx(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, result1)

	start2 := end1
	end2 := start2.Add(time.Hour)

	arg2 := BookingTxParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start2,
		End:          end2,
		PhoneBooking: phone,
	}
	result2, err := store.BookingTx(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, result2)
	require.NotEqual(t, result1.Booking.ID, result2.Booking.ID)
}
