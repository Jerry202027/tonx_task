package service

import (
	"tonx_task/database"
	"tonx_task/store"

	"gorm.io/gorm"
)

var (
	FlightBookingService = &FlightBookingServiceImpl{}
	FlightBookingStore   = &store.FlightBookingStore{}
	db                   *gorm.DB
)

// ServiceInit initializes the database and injects the Store into the ServiceImpl
func ServiceInit() {
	db = database.GetDB()

	// assigns the flightStore to FlightBookingStore
	*FlightBookingService = FlightBookingServiceImpl{
		flightStore: FlightBookingStore,
	}
}
