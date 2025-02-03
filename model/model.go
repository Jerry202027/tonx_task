package model

import "time"

type Flight struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FlightNo             string    `gorm:"size:20;not null" json:"flight_no"`
	Origin               string    `gorm:"size:3;not null" json:"origin"`
	Destination          string    `gorm:"size:3;not null" json:"destination"`
	FlightDate           time.Time `gorm:"not null" json:"flight_date"`
	DepartureTime        time.Time `gorm:"not null" json:"departure_time"`
	ArrivalTime          time.Time `gorm:"not null" json:"arrival_time"`
	Price                float64   `gorm:"not null" json:"price"`
	Capacity             int       `gorm:"not null" json:"capacity"`
	OverbookingThreshold int       `gorm:"not null;default:0" json:"overbooking_threshold"`
	Version              int       `gorm:"not null;default:1" json:"version"` // for optimistic lock
	RemainingSeats       int       `gorm:"not null" json:"remaining_seats"`   // initial value should be Capacity + OverbookingThreshold
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type Booking struct {
	ID            uint    `gorm:"primaryKey;autoIncrement"`
	FlightID      uint    `gorm:"index;not null"`
	Flight        Flight  `gorm:"foreignKey:FlightID;references:ID"`
	UserID        uint    `gorm:"not null"`
	BookingStatus string  `gorm:"size:20;not null"`
	FareClass     string  `gorm:"size:10;not null"`
	Price         float64 `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
