package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type SearchFlightRequest struct {
	Origin      string `form:"origin" binding:"required"`
	Destination string `form:"destination" binding:"required"`
	Date        string `form:"date" binding:"required"`
	ClassType   string `form:"class_type"`
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
}

type CreateFlightRequest struct {
	AirlineID     uint   `json:"airline_id" binding:"required"`
	OriginID      uint   `json:"origin_id" binding:"required"`
	DestinationID uint   `json:"destination_id" binding:"required"`
	DepartureTime string `json:"departure_time" binding:"required"`
	ArrivalTime   string `json:"arrival_time" binding:"required"`
	FlightNumber  string `json:"flight_number" binding:"required"`
}

type UpdateFlightRequest struct {
	AirlineID     uint   `json:"airline_id"`
	OriginID      uint   `json:"origin_id"`
	DestinationID uint   `json:"destination_id"`
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
	FlightNumber  string `json:"flight_number"`
}

type CreateTransactionRequest struct {
	FlightID    uint               `json:"flight_id" binding:"required"`
	ClassID     uint               `json:"class_id" binding:"required"`
	SeatNumbers []string           `json:"seat_numbers" binding:"required"`
	Passengers  []PassengerRequest `json:"passengers" binding:"required,dive"`
	PromoCode   *string            `json:"promo_code"`
}

type PassengerRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Nationality string `json:"nationality"`
	PassportNo  string `json:"passport_no" binding:"required"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type CreateAirlineRequest struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	LogoURL string `json:"logo_url"`
}

type UpdateAirlineRequest struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	LogoURL string `json:"logo_url"`
}

type CreateAirportRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
	City string `json:"city" binding:"required"`
}

type UpdateAirportRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

type CreatePromoRequest struct {
	Code     string  `json:"code" binding:"required"`
	Discount float64 `json:"discount" binding:"required,min=0,max=100"`
	IsActive bool    `json:"is_active"`
}

type UpdatePromoRequest struct {
	Code     string  `json:"code"`
	Discount float64 `json:"discount"`
	IsActive bool    `json:"is_active"`
}