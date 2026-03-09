package models

import (
	"time"

	"gorm.io/gorm"
)

type Airline struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Code      string         `gorm:"size:10;uniqueIndex;not null" json:"code"`
	LogoURL   string         `json:"logo_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Airport struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `gorm:"size:10;uniqueIndex;not null" json:"code"` // e.g., CGK
	Name string `gorm:"size:100;not null" json:"name"`
	City string `gorm:"size:100;not null" json:"city"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Flight struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	AirlineID     uint          `gorm:"not null;index" json:"airline_id"`
	Airline       Airline       `gorm:"foreignKey:AirlineID" json:"airline"`
	OriginID      uint          `gorm:"not null;index" json:"origin_id"`
	Origin        Airport       `gorm:"foreignKey:OriginID" json:"origin"`
	DestinationID uint          `gorm:"not null;index" json:"destination_id"`
	Destination   Airport       `gorm:"foreignKey:DestinationID" json:"destination"`
	DepartureTime time.Time     `gorm:"not null;index" json:"departure_time"`
	ArrivalTime   time.Time     `gorm:"not null" json:"arrival_time"`
	FlightNumber  string        `gorm:"size:20;not null" json:"flight_number"`
	TotalSeats    int           `gorm:"not null;default:0" json:"total_seats"`
	TotalRows     int           `gorm:"not null;default:0" json:"total_rows"`
	TotalColumns  int           `gorm:"not null;default:0" json:"total_columns"`
	FlightClasses []FlightClass `gorm:"foreignKey:FlightID" json:"classes"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type FlightClass struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	FlightID    uint        `gorm:"not null;index" json:"flight_id"`
	ClassType   string      `gorm:"size:20;not null;index" json:"class_type"`
	Price       float64     `gorm:"not null" json:"price"`
	Seats       []FlightSeat `gorm:"foreignKey:FlightClassID" json:"seats"`
	CreatedAt   time.Time   `json:"-"`
	UpdatedAt   time.Time   `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type FlightSeat struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	FlightClassID uint        `gorm:"not null;index" json:"flight_class_id"`
	SeatNumber    string      `gorm:"size:10;not null" json:"seat_number"`
	IsAvailable   bool        `gorm:"not null;index" json:"is_available"`
	FlightClass   FlightClass `gorm:"foreignKey:FlightClassID" json:"-"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
