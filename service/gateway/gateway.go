package gateway

import (
	"time"
	"tonx_task/model"

	"gorm.io/gorm"
)

type FlightBookingStore interface {
	AutoMigrate() error
	SearchFlights(origin, destination string, flightDate time.Time, offset, limit int) ([]model.Flight, error)
	CreateBooking(tx *gorm.DB, booking *model.Booking) error
	CountConfirmedBookings(tx *gorm.DB, flightID uint) (int, error)
}
