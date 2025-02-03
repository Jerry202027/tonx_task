package store

import (
	"fmt"
	"time"
	"tonx_task/database"
	"tonx_task/model"
	"tonx_task/service/gateway"

	"gorm.io/gorm"
)

type FlightBookingStore struct {
}

// check if FlightBookingStore implement gateway.FlightBookingStore
var _ gateway.FlightBookingStore = (*FlightBookingStore)(nil)

func (fs *FlightBookingStore) AutoMigrate() error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("db not initialized")
	}

	if err := db.AutoMigrate(&model.Flight{}, &model.Booking{}); err != nil {
		return err
	}

	return nil
}

// get flight data with the given filters (origin, destination, flightDate) + offset/limit
func (fs *FlightBookingStore) SearchFlights(origin, destination string, flightDate time.Time, offset, pageSize int) ([]model.Flight, error) {
	db := database.GetDB()
	var flights []model.Flight

	query := db.Model(&model.Flight{})
	if origin != "" {
		query = query.Where("origin = ?", origin)
	}

	if destination != "" {
		query = query.Where("destination = ?", destination)
	}

	if !flightDate.IsZero() {
		query = query.Where("DATE(flight_date) = DATE(?)", flightDate)
	}

	err := query.Offset(offset).Limit(pageSize).Find(&flights).Error
	if err != nil {
		fmt.Println("SearchFlights in store failed, err: ", err)
		return nil, err
	}

	return flights, nil
}

func (fs *FlightBookingStore) CreateBooking(tx *gorm.DB, booking *model.Booking) error {
	if err := tx.Create(booking).Error; err != nil {
		fmt.Println("CreateBooking in store failed, err: ", err)
		return err
	}

	return nil
}

func (fs *FlightBookingStore) CountConfirmedBookings(tx *gorm.DB, flightID uint) (int, error) {
	var count int64
	if err := tx.Model(&model.Booking{}).
		Where("flight_id = ? AND booking_status = 'CONFIRMED'", flightID).
		Count(&count).Error; err != nil {
		fmt.Println("CountConfirmedBookings in store failed, err: ", err)
		return 0, err
	}

	return int(count), nil
}
