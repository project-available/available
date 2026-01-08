package db

import (
	"context"
	"testing"
	"time"

	"github.com/project-available/available/utils"
	"github.com/stretchr/testify/require"
)

func createRandomBooking(t *testing.T) Booking {
	account := createRandomAccount(t)
	room := createRandomRoom(t)
	phoneBooking := utils.RandomPhone()
	start := time.Now().UTC()
	end := start.Add(time.Hour)

	arg := CreateBookingParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start,
		End:          end,
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

func TestCheckBookingOverlap(t *testing.T) {
	booking1 := createRandomBooking(t)

	// Overlap exact same time
	count, err := testQuery.CheckBookingOverlap(context.Background(), CheckBookingOverlapParams{
		RoomID: booking1.RoomID,
		End:    booking1.Start,
		Start:  booking1.End,
	})
	require.NoError(t, err)
	// Status of random booking might not be pending/confirmed, so we need to ensure it is.
	// But generated random booking doesn't set status? Default in DB probably pending or confirmed?
	// Checking the schema/queries would be good but based on generated code: CreateBooking accepts status? No, it returns status.
	// Let's assume default status is compatible or update acts on it.
	// Wait, CheckBookingOverlap filters by status IN ('pending', 'confirmed').
	// Let's verify what createRandomBooking makes. Db schema not visible but usually defaults to one of those.
	// If it fails we might need to update status to confirmed.

	// Let's update status to be sure
	updatedBooking, err := testQuery.UpdateBooking(context.Background(), UpdateBookingParams{
		ID:     booking1.ID,
		Status: "confirmed",
	})
	require.NoError(t, err)

	count, err = testQuery.CheckBookingOverlap(context.Background(), CheckBookingOverlapParams{
		RoomID: updatedBooking.RoomID,
		End:    updatedBooking.Start,
		Start:  updatedBooking.End,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), count)

	// No Overlap
	count, err = testQuery.CheckBookingOverlap(context.Background(), CheckBookingOverlapParams{
		RoomID: updatedBooking.RoomID,
		End:    updatedBooking.End.Add(time.Minute),
		Start:  updatedBooking.End.Add(time.Hour),
	})
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

func TestGetBookingsOnDate(t *testing.T) {
	account := createRandomAccount(t)
	room := createRandomRoom(t)

	start := time.Now().UTC()
	end := start.Add(2 * time.Hour)

	arg := CreateBookingParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start,
		End:          end,
		PhoneBooking: utils.RandomPhone(),
	}
	booking, err := testQuery.CreateBooking(context.Background(), arg)
	require.NoError(t, err)

	_, err = testQuery.UpdateBooking(context.Background(), UpdateBookingParams{
		ID:     booking.ID,
		Status: "confirmed",
	})
	require.NoError(t, err)

	dateStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	dateEnd := time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 999999999, start.Location())

	bookings, err := testQuery.GetBookingsOnDate(context.Background(), GetBookingsOnDateParams{
		RoomID:  room.ID,
		Start:   dateStart,
		Start_2: dateEnd,
	})
	require.NoError(t, err)
	require.NotEmpty(t, bookings)

	found := false
	for _, b := range bookings {
		if b.ID == booking.ID {
			found = true
			require.Equal(t, "confirmed", b.Status)
			break
		}
	}
	require.True(t, found)
}

func TestCheckBookingOverlap_PartialOverlap(t *testing.T) {
	account := createRandomAccount(t)
	room := createRandomRoom(t)

	start := time.Now().UTC()
	end := start.Add(2 * time.Hour)

	booking, err := testQuery.CreateBooking(context.Background(), CreateBookingParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start,
		End:          end,
		PhoneBooking: utils.RandomPhone(),
	})
	require.NoError(t, err)

	_, err = testQuery.UpdateBooking(context.Background(), UpdateBookingParams{
		ID:     booking.ID,
		Status: "confirmed",
	})
	require.NoError(t, err)

	overlapStart := start.Add(time.Hour)
	overlapEnd := end.Add(time.Hour)

	count, err := testQuery.CheckBookingOverlap(context.Background(), CheckBookingOverlapParams{
		RoomID: room.ID,
		End:    overlapStart,
		Start:  overlapEnd,
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), count)
}

func TestCheckBookingOverlap_DifferentRoom(t *testing.T) {
	account := createRandomAccount(t)
	room1 := createRandomRoom(t)
	room2 := createRandomRoom(t)

	start := time.Now().UTC()
	end := start.Add(time.Hour)

	booking, err := testQuery.CreateBooking(context.Background(), CreateBookingParams{
		AccountID:    account.ID,
		RoomID:       room1.ID,
		Start:        start,
		End:          end,
		PhoneBooking: utils.RandomPhone(),
	})
	require.NoError(t, err)

	_, err = testQuery.UpdateBooking(context.Background(), UpdateBookingParams{
		ID:     booking.ID,
		Status: "confirmed",
	})
	require.NoError(t, err)

	count, err := testQuery.CheckBookingOverlap(context.Background(), CheckBookingOverlapParams{
		RoomID: room2.ID,
		End:    start,
		Start:  end,
	})
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}

func TestCheckBookingOverlap_CancelledStatus(t *testing.T) {
	account := createRandomAccount(t)
	room := createRandomRoom(t)

	start := time.Now().UTC()
	end := start.Add(time.Hour)

	booking, err := testQuery.CreateBooking(context.Background(), CreateBookingParams{
		AccountID:    account.ID,
		RoomID:       room.ID,
		Start:        start,
		End:          end,
		PhoneBooking: utils.RandomPhone(),
	})
	require.NoError(t, err)

	_, err = testQuery.UpdateBooking(context.Background(), UpdateBookingParams{
		ID:     booking.ID,
		Status: "cancelled",
	})
	require.NoError(t, err)

	count, err := testQuery.CheckBookingOverlap(context.Background(), CheckBookingOverlapParams{
		RoomID: room.ID,
		End:    start,
		Start:  end,
	})
	require.NoError(t, err)
	require.Equal(t, int64(0), count)
}
