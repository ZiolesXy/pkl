package handler

import (
	"voca-plane/internal/domain/dto"
	"voca-plane/internal/domain/models"
	"voca-plane/internal/service"
	"voca-plane/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *service.AdminService
}

func NewAdminHandler(s *service.AdminService) *AdminHandler {
	return &AdminHandler{service: s}
}

func (h *AdminHandler) GetDashboard(c *gin.Context) {
	stats, err := h.service.GetDashboardStats(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "dashboard stats retrieved", stats)
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.service.GetAllUsers(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "users retrieved", users, meta)
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.UpdateUserRole(c.Request.Context(), uint(userID), req.Role)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user role updated", nil)
}

func (h *AdminHandler) GetTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	transactions, total, err := h.service.GetAllTransactions(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "transactions retrieved", transactions, meta)
}

func (h *AdminHandler) GetFlights(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	flights, total, err := h.service.GetAllFlights(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "flights retrieved", flights, meta)
}

func (h *AdminHandler) CreateFlight(c *gin.Context) {
	var req dto.CreateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	departureTime, _ := time.Parse(time.RFC3339, req.DepartureTime)
	arrivalTime, _ := time.Parse(time.RFC3339, req.ArrivalTime)

	flight := &models.Flight{
		AirlineID:     req.AirlineID,
		OriginID:      req.OriginID,
		DestinationID: req.DestinationID,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		FlightNumber:  req.FlightNumber,
	}

	err := h.service.CreateFlight(c.Request.Context(), flight)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "flight created", flight)
}

func (h *AdminHandler) UpdateFlight(c *gin.Context) {
	id := c.Param("id")
	flightID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid flight id")
		return
	}

	var req dto.UpdateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	flight := &models.Flight{ID: uint(flightID)}
	if req.AirlineID > 0 {
		flight.AirlineID = req.AirlineID
	}
	if req.OriginID > 0 {
		flight.OriginID = req.OriginID
	}
	if req.DestinationID > 0 {
		flight.DestinationID = req.DestinationID
	}
	if req.FlightNumber != "" {
		flight.FlightNumber = req.FlightNumber
	}

	err = h.service.UpdateFlight(c.Request.Context(), flight)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "flight updated", flight)
}

func (h *AdminHandler) DeleteFlight(c *gin.Context) {
	id := c.Param("id")
	flightID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid flight id")
		return
	}

	err = h.service.DeleteFlight(c.Request.Context(), uint(flightID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "flight deleted", nil)
}

func (h *AdminHandler) GetAirlines(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	airlines, total, err := h.service.GetAllAirlines(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "airlines retrieved", airlines, meta)
}

func (h *AdminHandler) CreateAirline(c *gin.Context) {
	var req dto.CreateAirlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	airline := &models.Airline{
		Name:    req.Name,
		Code:    req.Code,
		LogoURL: req.LogoURL,
	}

	err := h.service.CreateAirline(c.Request.Context(), airline)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "airline created", airline)
}

func (h *AdminHandler) UpdateAirline(c *gin.Context) {
	id := c.Param("id")
	airlineID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid airline id")
		return
	}

	var req dto.UpdateAirlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	airline := &models.Airline{ID: uint(airlineID)}
	if req.Name != "" {
		airline.Name = req.Name
	}
	if req.Code != "" {
		airline.Code = req.Code
	}
	if req.LogoURL != "" {
		airline.LogoURL = req.LogoURL
	}

	err = h.service.UpdateAirline(c.Request.Context(), airline)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "airline updated", airline)
}

func (h *AdminHandler) DeleteAirline(c *gin.Context) {
	id := c.Param("id")
	airlineID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid airline id")
		return
	}

	err = h.service.DeleteAirline(c.Request.Context(), uint(airlineID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "airline deleted", nil)
}

func (h *AdminHandler) GetAirports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	airports, total, err := h.service.GetAllAirports(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "airports retrieved", airports, meta)
}

func (h *AdminHandler) CreateAirport(c *gin.Context) {
	var req dto.CreateAirportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	airport := &models.Airport{
		Code: req.Code,
		Name: req.Name,
		City: req.City,
	}

	err := h.service.CreateAirport(c.Request.Context(), airport)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "airport created", airport)
}

func (h *AdminHandler) UpdateAirport(c *gin.Context) {
	id := c.Param("id")
	airportID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid airport id")
		return
	}

	var req dto.UpdateAirportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	airport := &models.Airport{ID: uint(airportID)}
	if req.Code != "" {
		airport.Code = req.Code
	}
	if req.Name != "" {
		airport.Name = req.Name
	}
	if req.City != "" {
		airport.City = req.City
	}

	err = h.service.UpdateAirport(c.Request.Context(), airport)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "airport updated", airport)
}

func (h *AdminHandler) DeleteAirport(c *gin.Context) {
	id := c.Param("id")
	airportID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid airport id")
		return
	}

	err = h.service.DeleteAirport(c.Request.Context(), uint(airportID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "airport deleted", nil)
}

func (h *AdminHandler) GetPromos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	promos, total, err := h.service.GetAllPromos(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	meta := gin.H{"total": total, "page": page, "limit": limit}
	response.SuccessWithMeta(c, http.StatusOK, "promos retrieved", promos, meta)
}

func (h *AdminHandler) CreatePromo(c *gin.Context) {
	var req dto.CreatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	promo := &models.PromoCode{
		Code:     req.Code,
		Discount: req.Discount,
		IsActive: req.IsActive,
	}

	err := h.service.CreatePromo(c.Request.Context(), promo)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "promo created", promo)
}

func (h *AdminHandler) UpdatePromo(c *gin.Context) {
	id := c.Param("id")
	promoID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid promo id")
		return
	}

	var req dto.UpdatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	promo := &models.PromoCode{ID: uint(promoID)}
	if req.Code != "" {
		promo.Code = req.Code
	}
	if req.Discount > 0 {
		promo.Discount = req.Discount
	}
	promo.IsActive = req.IsActive

	err = h.service.UpdatePromo(c.Request.Context(), promo)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "promo updated", promo)
}

func (h *AdminHandler) DeletePromo(c *gin.Context) {
	id := c.Param("id")
	promoID, err := strconv.Atoi(id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid promo id")
		return
	}

	err = h.service.DeletePromo(c.Request.Context(), uint(promoID))
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "promo deleted", nil)
}