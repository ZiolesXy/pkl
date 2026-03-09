package service

import (
	"context"
	"errors"
	"fmt"
	"voca-plane/internal/domain/dto"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"

	"gorm.io/gorm"
)

type AdminService struct {
	adminRepo   repository.AdminRepository
	userRepo    repository.UserRepository
	flightRepo  repository.FlightRepository
	airlineRepo repository.AirlineRepository
	airportRepo repository.AirportRepository
	promoRepo   repository.PromoRepository
	db          *gorm.DB
}

func NewAdminService(
	adminRepo repository.AdminRepository,
	userRepo repository.UserRepository,
	flightRepo repository.FlightRepository,
	airlineRepo repository.AirlineRepository,
	airportRepo repository.AirportRepository,
	promoRepo repository.PromoRepository,
	db *gorm.DB,
) *AdminService {
	return &AdminService{
		adminRepo:   adminRepo,
		userRepo:    userRepo,
		flightRepo:  flightRepo,
		airlineRepo: airlineRepo,
		airportRepo: airportRepo,
		promoRepo:   promoRepo,
		db:          db,
	}
}

func (s *AdminService) GetDashboardStats(ctx context.Context) (*repository.DashboardStats, error) {
	return s.adminRepo.GetDashboardStats(ctx)
}

func (s *AdminService) GetAllUsers(ctx context.Context, page, limit int) ([]models.User, int64, error) {
	return s.adminRepo.GetAllUsers(ctx, page, limit)
}

func (s *AdminService) UpdateUserRole(ctx context.Context, userID uint, role string) error {
	validRoles := map[string]bool{
		models.RoleUser:       true,
		models.RoleAdmin:      true,
		models.RoleSuperAdmin: true,
	}
	if !validRoles[role] {
		return errors.New("invalid role")
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.adminRepo.UpdateUserRole(ctx, tx, userID, role); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) GetAllTransactions(ctx context.Context, page, limit int) ([]models.Transaction, int64, error) {
	return s.adminRepo.GetAllTransactions(ctx, page, limit)
}

func (s *AdminService) GetAllFlights(ctx context.Context, page, limit int) ([]dto.FlightResponse, int64, error) {
	flights, total, err := s.flightRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var flightResponses []dto.FlightResponse
	for _, f := range flights {
		flightResponses = append(flightResponses, dto.ToFlightResponse(f))
	}

	return flightResponses, total, nil
}

func (s *AdminService) CreateFlight(ctx context.Context, flight *models.Flight, classCount int, classPrices []dto.ClassPriceRequest) (*dto.FlightResponse, error) {
	// Validate total_seats <= rows * columns
	maxCapacity := flight.TotalRows * flight.TotalColumns
	if flight.TotalSeats > maxCapacity {
		return nil, fmt.Errorf("total_seats (%d) exceeds rows × columns (%d × %d = %d)", flight.TotalSeats, flight.TotalRows, flight.TotalColumns, maxCapacity)
	}

	// Validate class_count
	if classCount < 1 || classCount > 3 {
		return nil, errors.New("class_count must be 1, 2, or 3")
	}

	// Validate class_prices count matches class_count
	if len(classPrices) != classCount {
		return nil, fmt.Errorf("class_prices count (%d) must match class_count (%d)", len(classPrices), classCount)
	}

	// Determine seat distribution by class
	type classAlloc struct {
		ClassType string
		Price     float64
		SeatCount int
	}

	var allocations []classAlloc
	totalSeats := flight.TotalSeats

	switch classCount {
	case 1:
		// Economy 100%
		allocations = append(allocations, classAlloc{
			ClassType: classPrices[0].ClassType,
			Price:     classPrices[0].Price,
			SeatCount: totalSeats,
		})
	case 2:
		// Business 30%, Economy 70%
		businessSeats := int(float64(totalSeats) * 0.30)
		economySeats := totalSeats - businessSeats
		allocations = append(allocations,
			classAlloc{ClassType: classPrices[0].ClassType, Price: classPrices[0].Price, SeatCount: businessSeats},
			classAlloc{ClassType: classPrices[1].ClassType, Price: classPrices[1].Price, SeatCount: economySeats},
		)
	case 3:
		// First 20%, Business 30%, Economy 50%
		firstSeats := int(float64(totalSeats) * 0.20)
		businessSeats := int(float64(totalSeats) * 0.30)
		economySeats := totalSeats - firstSeats - businessSeats
		allocations = append(allocations,
			classAlloc{ClassType: classPrices[0].ClassType, Price: classPrices[0].Price, SeatCount: firstSeats},
			classAlloc{ClassType: classPrices[1].ClassType, Price: classPrices[1].Price, SeatCount: businessSeats},
			classAlloc{ClassType: classPrices[2].ClassType, Price: classPrices[2].Price, SeatCount: economySeats},
		)
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the flight first
	if err := tx.WithContext(ctx).Create(flight).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Generate seats per class
	seatIndex := 0 // global seat counter across all classes
	columns := flight.TotalColumns

	for _, alloc := range allocations {
		fClass := models.FlightClass{
			FlightID:  flight.ID,
			ClassType: alloc.ClassType,
			Price:     alloc.Price,
		}
		if err := tx.WithContext(ctx).Create(&fClass).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		for i := 0; i < alloc.SeatCount; i++ {
			row := seatIndex / columns       // 0-based row index
			col := (seatIndex % columns) + 1 // 1-based column
			rowLetter := string(rune('A' + row))
			seatNumber := fmt.Sprintf("%s%d", rowLetter, col)

			seat := models.FlightSeat{
				FlightClassID: fClass.ID,
				SeatNumber:    seatNumber,
				IsAvailable:   true,
			}
			if err := tx.WithContext(ctx).Create(&seat).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			seatIndex++
		}
	}

	if err := tx.WithContext(ctx).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses.Seats").
		First(flight, flight.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	response := dto.ToFlightResponse(*flight)
	return &response, nil
}

func (s *AdminService) UpdateFlight(ctx context.Context, flight *models.Flight) (*dto.FlightResponse, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Fetch old flight with classes and seats to compare
	var oldFlight models.Flight
	if err := tx.WithContext(ctx).
		Preload("FlightClasses.Seats").
		First(&oldFlight, flight.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Validate new seat config
	if flight.TotalSeats > 0 && flight.TotalRows > 0 && flight.TotalColumns > 0 {
		maxCapacity := flight.TotalRows * flight.TotalColumns
		if flight.TotalSeats > maxCapacity {
			tx.Rollback()
			return nil, fmt.Errorf("total_seats (%d) exceeds rows × columns (%d × %d = %d)", flight.TotalSeats, flight.TotalRows, flight.TotalColumns, maxCapacity)
		}
	}

	// Update ONLY flight table columns using map — avoids GORM re-saving associations
	if err := tx.WithContext(ctx).Model(&models.Flight{}).Where("id = ?", flight.ID).Updates(map[string]interface{}{
		"airline_id":      flight.AirlineID,
		"origin_id":       flight.OriginID,
		"destination_id":  flight.DestinationID,
		"departure_time":  flight.DepartureTime,
		"arrival_time":    flight.ArrivalTime,
		"flight_number":   flight.FlightNumber,
		"total_seats":     flight.TotalSeats,
		"total_rows":      flight.TotalRows,
		"total_columns":   flight.TotalColumns,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Recalculate seats if total_seats, total_rows, or total_columns changed
	seatsChanged := oldFlight.TotalSeats != flight.TotalSeats ||
		oldFlight.TotalRows != flight.TotalRows ||
		oldFlight.TotalColumns != flight.TotalColumns
	classCount := len(oldFlight.FlightClasses)

	if seatsChanged && classCount > 0 {
		newTotal := flight.TotalSeats

		// Calculate new seat distribution per class
		type classTarget struct {
			ClassID   uint
			SeatCount int
		}
		var targets []classTarget

		switch classCount {
		case 1:
			targets = append(targets, classTarget{oldFlight.FlightClasses[0].ID, newTotal})
		case 2:
			biz := int(float64(newTotal) * 0.30)
			eco := newTotal - biz
			targets = append(targets,
				classTarget{oldFlight.FlightClasses[0].ID, biz},
				classTarget{oldFlight.FlightClasses[1].ID, eco},
			)
		case 3:
			first := int(float64(newTotal) * 0.20)
			biz := int(float64(newTotal) * 0.30)
			eco := newTotal - first - biz
			targets = append(targets,
				classTarget{oldFlight.FlightClasses[0].ID, first},
				classTarget{oldFlight.FlightClasses[1].ID, biz},
				classTarget{oldFlight.FlightClasses[2].ID, eco},
			)
		}

		// Hard delete all existing seats (Unscoped to bypass soft delete)
		for _, fc := range oldFlight.FlightClasses {
			if err := tx.WithContext(ctx).Unscoped().Where("flight_class_id = ?", fc.ID).Delete(&models.FlightSeat{}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Regenerate seats
		columns := flight.TotalColumns
		seatIndex := 0
		for _, t := range targets {
			for i := 0; i < t.SeatCount; i++ {
				row := seatIndex / columns
				col := (seatIndex % columns) + 1
				rowLetter := string(rune('A' + row))
				seatNumber := fmt.Sprintf("%s%d", rowLetter, col)

				seat := models.FlightSeat{
					FlightClassID: t.ClassID,
					SeatNumber:    seatNumber,
					IsAvailable:   true,
				}
				if err := tx.WithContext(ctx).Create(&seat).Error; err != nil {
					tx.Rollback()
					return nil, err
				}
				seatIndex++
			}
		}
	}

	// Reload with preloads
	if err := tx.WithContext(ctx).
		Preload("Airline").
		Preload("Origin").
		Preload("Destination").
		Preload("FlightClasses.Seats").
		First(flight, flight.ID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	response := dto.ToFlightResponse(*flight)
	return &response, nil
}

func (s *AdminService) DeleteFlight(ctx context.Context, id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.flightRepo.Delete(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) GetAllAirlines(ctx context.Context, page, limit int) ([]models.Airline, int64, error) {
	return s.airlineRepo.GetAll(ctx, page, limit)
}

func (s *AdminService) CreateAirline(ctx context.Context, airline *models.Airline) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airlineRepo.Create(ctx, tx, airline); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) UpdateAirline(ctx context.Context, airline *models.Airline) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airlineRepo.Update(ctx, tx, airline); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) DeleteAirline(ctx context.Context, id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airlineRepo.Delete(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) GetAllAirports(ctx context.Context, page, limit int) ([]models.Airport, int64, error) {
	return s.airportRepo.GetAll(ctx, page, limit)
}

func (s *AdminService) CreateAirport(ctx context.Context, airport *models.Airport) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airportRepo.Create(ctx, tx, airport); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) UpdateAirport(ctx context.Context, airport *models.Airport) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airportRepo.Update(ctx, tx, airport); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) DeleteAirport(ctx context.Context, id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.airportRepo.Delete(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) GetAllPromos(ctx context.Context, page, limit int) ([]models.PromoCode, int64, error) {
	return s.promoRepo.GetAll(ctx, page, limit)
}

func (s *AdminService) CreatePromo(ctx context.Context, promo *models.PromoCode) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.promoRepo.Create(ctx, tx, promo); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) UpdatePromo(ctx context.Context, promo *models.PromoCode) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.promoRepo.Update(ctx, tx, promo); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) DeletePromo(ctx context.Context, id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.promoRepo.Delete(ctx, tx, id); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) GetFlightByID(ctx context.Context, id uint) (*dto.FlightResponse, error) {
	flight, err := s.flightRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := dto.ToFlightResponse(*flight)
	return &response, nil
}

func (s *AdminService) GetFlightModelByID(ctx context.Context, id uint) (*models.Flight, error) {
	return s.flightRepo.GetByID(ctx, id)
}

func (s *AdminService) GetAirlineByID(ctx context.Context, id uint) (*models.Airline, error) {
	return s.airlineRepo.GetByID(ctx, id)
}

func (s *AdminService) GetAirportByID(ctx context.Context, id uint) (*models.Airport, error) {
	return s.airportRepo.GetByID(ctx, id)
}

func (s *AdminService) GetPromoByID(ctx context.Context, id uint) (*models.PromoCode, error) {
	return s.promoRepo.GetByID(ctx, id)
}
