package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"focastore/response"
	"focastore/seeders"
)

type SeedHandler struct {
	db *gorm.DB
}

func NewSeedHandler(db *gorm.DB) *SeedHandler {
	return &SeedHandler{db: db}
}

func (h *SeedHandler) Seed(c *gin.Context) {
	if err := seeders.Seed(h.db); err != nil {
		response.Error(c, http.StatusInternalServerError, "seed failed")
		return
	}
	response.Success(c, "seeded", gin.H{})
}
