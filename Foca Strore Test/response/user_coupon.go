package response

import (
	"time"
	"voca-store/models"
)

type UserCouponResponse struct {
	ID         uint       `json:"id"`
	CouponCode string     `json:"coupon_code"`
	CouponType string     `json:"coupon_type"`
	Value      float64    `json:"value"`
	UsedAt     *time.Time `json:"used_at"`
	ClaimedAt  time.Time  `json:"claimed_at"`
}

type UserCouponListResponse struct {
	Entries []UserCouponResponse `json:"entries"`
}

func BuildUserCouponResponse(uc models.UserCoupon) UserCouponResponse {
	return UserCouponResponse{
		ID:         uc.ID,
		CouponCode: uc.Coupon.Code,
		CouponType: uc.Coupon.Type,
		Value:      uc.Coupon.Value,
		UsedAt:     uc.UsedAt,
		ClaimedAt:  uc.CreatedAt,
	}
}

func BuildUserCouponListResponse(ucs []models.UserCoupon) UserCouponListResponse {
	entries := []UserCouponResponse{}
	for _, uc := range ucs {
		entries = append(entries, BuildUserCouponResponse(uc))
	}
	return UserCouponListResponse{Entries: entries}
}
