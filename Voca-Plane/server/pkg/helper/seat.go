package helper

import (
	"errors"
	"fmt"
	"strings"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/domain/models"
)

const (
	FirstClassRatio    = 0.20
	BusinessClassRatio = 0.30
)

type ClassAlloc struct {
	ClassType string
	Price     float64
	SeatCount int
}

func ValidateFlightInput(flight *models.Flight, classCount int, classPrice []request.ClassPriceRequest) error {
	maxCapacity := flight.TotalRows * flight.TotalColumns

	if flight.TotalSeats > maxCapacity {
		return fmt.Errorf("total_seats (%d) exceeds row x columns (%d)", flight.TotalSeats, maxCapacity)
	}

	if classCount < 1 || classCount > 3 {
		return errors.New("class_count must be 1, 2, ro 3")
	}

	if len(classPrice) != classCount {
		return fmt.Errorf("class_pries count (%d) must match class_count (%d)", len(classPrice), classCount)
	}

	return nil
}

func MapClassPrices(classPrices []request.ClassPriceRequest) map[string]float64 {
	classMap := make(map[string]float64)

	for _, cp := range classPrices {
		classType := strings.ToLower(cp.ClassType)
		classMap[classType] = cp.Price
	}

	return classMap
}

func CalculateSeatAllocation(totalSeats int, classMap map[string]float64, classCount int) []ClassAlloc {
	var allocations []ClassAlloc

	switch classCount {

	case 1:
		allocations = append(allocations, 
		ClassAlloc{"economy", classMap["economy"], totalSeats})
	
	case 2:
		businessSeats := int(float64(totalSeats) * BusinessClassRatio)
		economySeats := totalSeats - businessSeats

		allocations = append(allocations, 
			ClassAlloc{"business", classMap["business"], businessSeats},
			ClassAlloc{"economy", classMap["economy"], economySeats},
		)

	case 3:
		firstSeats := int(float64(totalSeats) * FirstClassRatio)
		businessSeats := int(float64(totalSeats) * BusinessClassRatio)
		economySeats := totalSeats - firstSeats - businessSeats

		allocations = append(allocations,
			ClassAlloc{"first", classMap["first"], firstSeats},
			ClassAlloc{"business", classMap["business"], businessSeats},
			ClassAlloc{"economy", classMap["economy"], economySeats},
		)
	}

	return allocations
}

func GenerateSeats(flightClassID uint, count int, startIndex int, columns int) []models.FlightSeat {
	seats := make([]models.FlightSeat, 0, count)

	for i := 0; i < count; i++ {
		index := startIndex + i

		row := index/columns
		col := (index % columns) + 1

		rowLetter := string(rune('A' + row))
		seatNumber := fmt.Sprintf("%s%d", rowLetter, col)

		seats = append(seats, models.FlightSeat{
			FlightClassID: flightClassID,
			SeatNumber: seatNumber,
			IsAvailable: true,
		})
	}

	return seats
}