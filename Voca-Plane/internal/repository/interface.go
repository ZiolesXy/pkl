package repository

import (
	"context"
	"voca-plane/internal/domain/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
}

type FlightRepository interface {
	Search(ctx context.Context, origin, destination, date, classType string, page, limit int) ([]models.Flight, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Flight, error)
	GetClassByID(ctx context.Context, id uint) (*models.FlightClass, error)
	GetSeat(ctx context.Context, classID uint, seatNumber string) (*models.FlightSeat, error)
	UpdateSeatAvailability(ctx context.Context, tx *gorm.DB, seatID uint, available bool) error
	GetAll(ctx context.Context, page, limit int) ([]models.Flight, int64, error)
	Create(ctx context.Context, tx *gorm.DB, flight *models.Flight) error
	Update(ctx context.Context, tx *gorm.DB, flight *models.Flight) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
}

type AirlineRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]models.Airline, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Airline, error)
	Create(ctx context.Context, tx *gorm.DB, airline *models.Airline) error
	Update(ctx context.Context, tx *gorm.DB, airline *models.Airline) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
}

type AirportRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]models.Airport, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Airport, error)
	Create(ctx context.Context, tx *gorm.DB, airport *models.Airport) error
	Update(ctx context.Context, tx *gorm.DB, airport *models.Airport) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, t *models.Transaction) error
	GetByCode(ctx context.Context, code string) (*models.Transaction, error)
	GetByUserID(ctx context.Context, userID uint, page, limit int) ([]models.Transaction, int64, error)
	UpdatePaymentStatus(ctx context.Context, tx *gorm.DB, id uint, status string) error
	Delete(ctx context.Context, tx *gorm.DB, code string) error
}

type PromoRepository interface {
	GetByID(ctx context.Context, id uint) (*models.PromoCode, error)
	GetByCode(ctx context.Context, code string) (*models.PromoCode, error)
	GetAll(ctx context.Context, page, limit int) ([]models.PromoCode, int64, error)
	Create(ctx context.Context, tx *gorm.DB, promo *models.PromoCode) error
	Update(ctx context.Context, tx *gorm.DB, promo *models.PromoCode) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
}

type DashboardStats struct {
	TotalUsers        int64   `json:"total_users"`
	TotalFlights      int64   `json:"total_flights"`
	TotalTransactions int64   `json:"total_transactions"`
	TotalRevenue      float64 `json:"total_revenue"`
	PendingPayments   int64   `json:"pending_payments"`
	CompletedBookings int64   `json:"completed_bookings"`
}

type AdminRepository interface {
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
	GetAllUsers(ctx context.Context, page, limit int) ([]models.User, int64, error)
	UpdateUserRole(ctx context.Context, tx *gorm.DB, userID uint, role string) error
	GetAllTransactions(ctx context.Context, page, limit int) ([]models.Transaction, int64, error)
}

type SystemRepository interface {
	ResetDatabase(ctx context.Context) error
}