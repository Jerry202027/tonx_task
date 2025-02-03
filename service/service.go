package service

import (
	"fmt"
	"time"
	"tonx_task/cache"
	"tonx_task/database"
	"tonx_task/model"
	"tonx_task/service/gateway"

	"gorm.io/gorm"
)

// Accesses the database through flightStore
type FlightBookingServiceImpl struct {
	flightStore gateway.FlightBookingStore
}

func (s *FlightBookingServiceImpl) AutoMigrate() error {
	return s.flightStore.AutoMigrate()
}

func (s *FlightBookingServiceImpl) SearchFlights(origin, destination string, flightDate time.Time, page, pageSize int) ([]model.Flight, error) {
	offset := (page - 1) * pageSize
	return s.flightStore.SearchFlights(origin, destination, flightDate, offset, pageSize)
}

func (s *FlightBookingServiceImpl) CreateBooking(flightID, userID uint, fareClass string, price float64) error {
	db := database.GetDB()
	if db == nil {
		return fmt.Errorf("db not initialized")
	}

	// 1. decrement the remaining seats using Redis + Lua script
	newRemaining, err := cache.DecrementRemainingSeats(flightID)
	if err != nil {
		return fmt.Errorf("failed to decrement remaining seats in Redis: %v", err)
	}

	// no seats available (including overbooking limit)
	if newRemaining < 0 {
		return fmt.Errorf("no seats available (including overbooking limit)")
	}

	// 2. db transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// 2.1 get flights data
		var flight model.Flight
		if err := tx.First(&flight, flightID).Error; err != nil {
			fmt.Println("get flights data in CreateBooking service failed, err: ", err)
			return err
		}

		if flight.ID == 0 {
			return fmt.Errorf("flight not found")
		}

		// 2.2 Check for overbooking
		confirmedCount, err := s.flightStore.CountConfirmedBookings(tx, flight.ID)
		if err != nil {
			return err
		}

		// maxAllowed = Capacity + OverbookingThreshold
		maxAllowed := flight.Capacity + flight.OverbookingThreshold
		if confirmedCount >= maxAllowed {
			return fmt.Errorf("no seats available (including overbooking limit)")
		}

		// 2.3 Use optimistic locking to update flights
		// save the current version
		oldVersion := flight.Version
		// update the flight's RemainingSeats; add version by 1
		result := tx.Model(&model.Flight{}).
			Where("id = ? AND version = ?", flight.ID, oldVersion).
			Updates(map[string]interface{}{
				"remaining_seats": newRemaining,
				"version":         oldVersion + 1,
			})

		if result.Error != nil {
			return result.Error
		}

		// a version conflict
		if result.RowsAffected == 0 {
			return fmt.Errorf("optimistic lock conflict, please retry")
		}

		// 3. Create a record
		newBooking := model.Booking{
			FlightID:      flight.ID,
			UserID:        userID,
			BookingStatus: "CONFIRMED",
			FareClass:     fareClass,
			Price:         price,
		}

		return s.flightStore.CreateBooking(tx, &newBooking)
	})
}
