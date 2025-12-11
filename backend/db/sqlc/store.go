package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Store interface {
	Querier
	BookingTx(ctx context.Context, arg BookingTxParams) (BookingTxResult, error)
}
type SQLStore struct {
	*Queries
	db *sql.DB
}
type SQLStoreq1 struct {
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// execTx executes a function within a database transaction
func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// The new function accept dbtx interface so sql.db or sql.tx are ok
	q := New(tx) // this is the query run within transaction

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type BookingTxParams struct {
	AccountID    int64     `json:"account_id"`
	RoomID       int64     `json:"room_id"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	PhoneBooking string    `json:"phone_booking"`
}

type BookingTxResult struct {
	Booking Booking `json:"booking"`
}

func (s *SQLStore) BookingTx(ctx context.Context, arg BookingTxParams) (BookingTxResult, error) {
	var result BookingTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		// Check for overlapping bookings on the same room
		overlap, err := s.CheckBookingOverlap(ctx, CheckBookingOverlapParams{
			RoomID: arg.RoomID,
			Start:  arg.Start,
			End:    arg.End,
		})
		if err != nil {
			return err
		}
		if overlap > 0 {
			return errors.New("booking time overlaps with existing bookings on this room")
		}

		// Create the booking
		arg1 := CreateBookingParams{
			AccountID:    arg.AccountID,
			RoomID:       arg.RoomID,
			Start:        arg.Start,
			End:          arg.End,
			PhoneBooking: arg.PhoneBooking,
		}
		booking, err := s.CreateBooking(ctx, arg1)
		if err != nil {
			return err
		}

		// Update status to confirmed
		arg2 := UpdateBookingParams{
			ID:     booking.ID,
			Status: "confirmed",
		}
		result.Booking, err = s.UpdateBooking(ctx, arg2)
		if err != nil {
			return err
		}

		return err
	})

	fmt.Println(result)

	return result, err
}
