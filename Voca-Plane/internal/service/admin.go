package service

import (
	"context"
	"errors"
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

func (s *AdminService) GetAllFlights(ctx context.Context, page, limit int) ([]models.Flight, int64, error) {
	return s.flightRepo.GetAll(ctx, page, limit)
}

func (s *AdminService) CreateFlight(ctx context.Context, flight *models.Flight) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.flightRepo.Create(ctx, tx, flight); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AdminService) UpdateFlight(ctx context.Context, flight *models.Flight) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := s.flightRepo.Update(ctx, tx, flight); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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

func (s *AdminService) GetFlightByID(ctx context.Context, id uint) (*models.Flight, error) {
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
