package service

import (
	"context"
	"errors"
	"fmt"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/domain/dto/response"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"
	"voca-plane/pkg/helper"

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

func (s *AdminService) GetAllFlights(ctx context.Context, page, limit int) ([]response.FlightResponse, int64, error) {
	flights, total, err := s.flightRepo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var flightResponses []response.FlightResponse
	for _, f := range flights {
		flightResponses = append(flightResponses, response.ToFlightResponse(f))
	}

	return flightResponses, total, nil
}

func (s *AdminService) CreateFlight(ctx context.Context, flight *models.Flight, classCount int, classPrices []request.ClassPriceRequest) (*response.FlightResponse, error) {
	if err := helper.ValidateFlightInput(flight, classCount, classPrices); err != nil {
		return nil, err
	}

	classMap := helper.MapClassPrices(classPrices)

	allocations := helper.CalculateSeatAllocation(
		flight.TotalSeats,
		classMap,
		classCount,
	)

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.flightRepo.Create(ctx, tx, flight); err != nil {
		tx.Rollback()
		return nil, err
	}

	seatIndex := 0
	columns := flight.TotalColumns

	for _, alloc := range allocations {

		fClass := &models.FlightClass{
			FlightID:  flight.ID,
			ClassType: alloc.ClassType,
			Price:     alloc.Price,
		}

		if err := s.flightRepo.CreateClass(ctx, tx, fClass); err != nil {
			tx.Rollback()
			return nil, err
		}

		seats := helper.GenerateSeats(
			fClass.ID,
			alloc.SeatCount,
			seatIndex,
			columns,
		)

		if err := s.flightRepo.BulkCreateSeats(ctx, tx, seats); err != nil {
			tx.Rollback()
			return nil, err
		}

		seatIndex += alloc.SeatCount
	}

	flight, err := s.flightRepo.GetFlightWithRelations(ctx, tx, flight.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToFlightResponse(*flight)

	return &res, nil
}

func (s *AdminService) UpdateFlight(ctx context.Context, flight *models.Flight) (*response.FlightResponse, error) {

	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1️⃣ Get existing flight
	oldFlight, err := s.flightRepo.GetFlightWithClasses(ctx, tx, flight.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2️⃣ Validate seat capacity
	maxCapacity := flight.TotalRows * flight.TotalColumns
	if flight.TotalSeats > maxCapacity {
		tx.Rollback()
		return nil, fmt.Errorf(
			"total_seats (%d) exceeds rows × columns (%d)",
			flight.TotalSeats,
			maxCapacity,
		)
	}

	// 3️⃣ Update flight
	if err := s.flightRepo.Update(ctx, tx, flight); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 4️⃣ Detect seat layout change
	seatsChanged :=
		oldFlight.TotalSeats != flight.TotalSeats ||
			oldFlight.TotalRows != flight.TotalRows ||
			oldFlight.TotalColumns != flight.TotalColumns

	if seatsChanged && len(oldFlight.FlightClasses) > 0 {

		classCount := len(oldFlight.FlightClasses)

		classMap := map[string]float64{}
		classIDs := make([]uint, 0)

		for _, fc := range oldFlight.FlightClasses {
			classMap[fc.ClassType] = fc.Price
			classIDs = append(classIDs, fc.ID)
		}

		allocations := helper.CalculateSeatAllocation(
			flight.TotalSeats,
			classMap,
			classCount,
		)

		// 5️⃣ Delete old seats
		if err := s.flightRepo.DeleteSeatsByClassIDs(ctx, tx, classIDs); err != nil {
			tx.Rollback()
			return nil, err
		}

		columns := flight.TotalColumns
		seatIndex := 0

		// 6️⃣ Regenerate seats
		for i, alloc := range allocations {

			seats := helper.GenerateSeats(
				oldFlight.FlightClasses[i].ID,
				alloc.SeatCount,
				seatIndex,
				columns,
			)

			if err := s.flightRepo.BulkCreateSeats(ctx, tx, seats); err != nil {
				tx.Rollback()
				return nil, err
			}

			seatIndex += alloc.SeatCount
		}
	}

	// 7️⃣ Reload flight
	flight, err = s.flightRepo.GetFlightWithRelations(ctx, tx, flight.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	res := response.ToFlightResponse(*flight)

	return &res, nil
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

func (s *AdminService) GetFlightByID(ctx context.Context, id uint) (*response.FlightResponse, error) {
	flight, err := s.flightRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := response.ToFlightResponse(*flight)
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
