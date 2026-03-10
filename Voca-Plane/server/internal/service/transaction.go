package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"voca-plane/internal/domain/dto/request"
	"voca-plane/internal/domain/dto/response"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/repository"
	"voca-plane/pkg/helper"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	txRepo     repository.TransactionRepository
	flightRepo repository.FlightRepository
	promoRepo  repository.PromoRepository
	db         *gorm.DB
	midtrans   *helper.MidTransClient
}

func NewTransactionService(txRepo repository.TransactionRepository, flightRepo repository.FlightRepository, promoRepo repository.PromoRepository, db *gorm.DB, midtrans *helper.MidTransClient) *TransactionService {
	return &TransactionService{
		txRepo:     txRepo,
		flightRepo: flightRepo,
		promoRepo:  promoRepo,
		db:         db,
		midtrans:   midtrans,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, userID uint, req request.CreateTransactionRequest) (*response.TransactionResponse, error) {
	flightClass, err := s.flightRepo.GetClassByID(ctx, req.ClassID)
	if err != nil || flightClass.FlightID != req.FlightID {
		return nil, errors.New("invalid flight class")
	}

	lockDuration := 10 * time.Minute
	lockedSeats := make([]*models.FlightSeat, 0)

	seatMap := map[string]bool{}
	for _, p := range req.Passengers {
		if seatMap[p.SeatNumber] {
			return nil, fmt.Errorf("duplicated seat %s", p.SeatNumber)
		}

		seatMap[p.SeatNumber] = true
	}

	for _, p := range req.Passengers {
		seat, err := s.flightRepo.GetSeat(ctx, req.ClassID, p.SeatNumber)
		if err != nil {
			return nil, fmt.Errorf("seat %s not found", p.SeatNumber)
		}

		if !seat.IsAvailable {
			return nil, fmt.Errorf("seat %s is already booked", p.SeatNumber)
		}

		lockedSeats = append(lockedSeats, seat)
	}

	subtotal := flightClass.Price * float64(len(req.Passengers))
	discount := 0.0

	if req.PromoCode != nil {
		promo, err := s.promoRepo.GetByCode(ctx, *req.PromoCode)
		if err == nil && promo.IsActive {
			discount = subtotal * (promo.Discount / 100)
		}
	}

	totalPrice := subtotal - discount

	var transaction models.Transaction
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, seat := range lockedSeats {
			if err := s.flightRepo.UpdateSeatAvailability(ctx, tx, seat.ID, false); err != nil {
				return err
			}
		}

		code := uuid.New().String()
		transaction = models.Transaction{
			Code:          code,
			UserID:        userID,
			FlightID:      req.FlightID,
			FlightClassID: req.ClassID,
			TotalPrice:    totalPrice,
			PaymentStatus: "PENDING",
			PromoCode:     req.PromoCode,
			Discount:      discount,
			ExpiresAt:     time.Now().Add(lockDuration),
		}

		if err := s.txRepo.Create(ctx, tx, &transaction); err != nil {
			return err
		}

		passengers := make([]models.TransactionPassenger, len(req.Passengers))
		for i, p := range req.Passengers {
			passengers[i] = models.TransactionPassenger{
				TransactionID: transaction.ID,
				FullName:      p.FullName,
				Nationality:   p.Nationality,
				PassportNo:    p.PassportNo,
				SeatNumber:    p.SeatNumber,
				FlightClassID: req.ClassID,
			}
		}
		if err := s.txRepo.CreatePassengers(ctx, tx, passengers); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if transaction.TotalPrice <= 0 {
		transaction.PaymentStatus = "PAID"
		transactionRes, err := s.txRepo.GetByCode(ctx, transaction.Code)
		if err != nil {
			return nil, err
		}
		resDto := response.ToTransactionResponse(*transactionRes)
		return &resDto, nil
	}

	res, err := s.midtrans.CreatePayment(
		transaction.Code,
		transaction.TotalPrice,
	)

	if err != nil {
		return nil, err
	}

	transaction.PaymentURL = res.RedirectURL

	if err := s.txRepo.UpdatePaymentURL(ctx, transaction.Code, res.RedirectURL); err != nil {
		return nil, err
	}

	transactionRes, err := s.txRepo.GetByCode(ctx, transaction.Code)
	if err != nil {
		return nil, err
	}

	resDto := response.ToTransactionResponse(*transactionRes)
	return &resDto, nil
}

func (s *TransactionService) PayTransaction(ctx context.Context, code string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction, err := s.txRepo.GetByCode(ctx, code)
	if err != nil {
		tx.Rollback()
		return errors.New("transaction not found")
	}

	if transaction.PaymentStatus == "PAID" {
		tx.Rollback()
		return errors.New("already paid")
	}

	if err := s.txRepo.UpdatePaymentStatus(ctx, tx, transaction.ID, "PAID"); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *TransactionService) GetUserTransactions(ctx context.Context, userID uint, page, limit int) ([]response.TransactionResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	transactions, total, err := s.txRepo.GetByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var res []response.TransactionResponse
	for _, t := range transactions {
		res = append(res, response.ToTransactionResponse(t))
	}

	return res, total, nil
}

func (s *TransactionService) GetTransactionByCode(ctx context.Context, code string) (*response.TransactionResponse, error) {
	transaction, err := s.txRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	res := response.ToTransactionResponse(*transaction)
	return &res, nil
}

func (s *TransactionService) CancelTransaction(ctx context.Context, userID uint, code string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	transaction, err := s.txRepo.GetByCode(ctx, code)
	if err != nil {
		tx.Rollback()
		return errors.New("transaction not found")
	}

	if transaction.UserID != userID {
		tx.Rollback()
		return errors.New("unauthorized")
	}

	if transaction.PaymentStatus == "PAID" {
		tx.Rollback()
		return errors.New("cannot cancel paid transaction")
	}

	if err := s.txRepo.Delete(ctx, tx, code); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
