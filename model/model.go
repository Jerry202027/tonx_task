package model

import "time"

type Flight struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement"`
	FlightNo             string    `gorm:"size:20;not null"`
	Origin               string    `gorm:"size:3;not null"`
	Destination          string    `gorm:"size:3;not null"`
	FlightDate           time.Time `gorm:"not null"`
	DepartureTime        time.Time `gorm:"not null"`
	ArrivalTime          time.Time `gorm:"not null"`
	Price                float64   `gorm:"not null"`
	Capacity             int       `gorm:"not null"`
	OverbookingThreshold int       `gorm:"not null;default:0"`
	Version              int       `gorm:"not null;default:1" json:"version"` // for optimistic lock
	RemainingSeats       int       `gorm:"not null" json:"remaining_seats"`   // initial value: Capacity + OverbookingThreshold
	CreatedAt            time.Time
	UpdatedAt            time.Time
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
