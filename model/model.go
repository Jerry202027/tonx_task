package model

import "time"

type Flight struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FlightNo             string    `gorm:"size:20;not null;index:idx_flight_no" json:"flight_no"`
	Origin               string    `gorm:"size:3;not null;index:idx_origin_destination_flightdate" json:"origin"`
	Destination          string    `gorm:"size:3;not null;index:idx_origin_destination_flightdate" json:"destination"`
	FlightDate           time.Time `gorm:"not null;index:idx_origin_destination_flightdate" json:"flight_date"`
	DepartureTime        time.Time `gorm:"not null" json:"departure_time"`
	ArrivalTime          time.Time `gorm:"not null" json:"arrival_time"`
	Price                float64   `gorm:"not null" json:"price"`
	Capacity             int       `gorm:"not null" json:"capacity"`
	OverbookingThreshold int       `gorm:"not null;default:0" json:"overbooking_threshold"`
	Version              int       `gorm:"not null;default:1" json:"version"` // For optimistic lock
	RemainingSeats       int       `gorm:"not null" json:"remaining_seats"`   // Initial value: Capacity + OverbookingThreshold
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Booking represents a booking record.
type Booking struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FlightID      uint      `gorm:"not null;index:idx_flight_id" json:"flight_id"`
	Flight        Flight    `gorm:"foreignKey:FlightID;references:ID" json:"-"`
	UserID        uint      `gorm:"not null" json:"user_id"`
	BookingStatus string    `gorm:"size:20;not null" json:"booking_status"` // e.g. CONFIRMED, CANCELLED...
	FareClass     string    `gorm:"size:10;not null" json:"fare_class"`
	Price         float64   `gorm:"not null" json:"price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
