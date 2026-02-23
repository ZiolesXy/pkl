package handlers

import (
	"errors"
	"os"
	"strings"

	// "fmt"
	"net/http"
	// "os"
	// "strings"
	// "voca-store/helper"
	"voca-store/helper"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exist := c.Get("user_id")
		if !exist {
			response.ErrorResponse(c, http.StatusUnauthorized, "user not authenticated")
			return
		}

		var user models.User
		if err := db.Preload("Role").First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(c, http.StatusNotFound, "user not found")
	
			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch user")
			}
			return 
		}

		profileResp := response.BuildUserProfileResponse(
			user.ID,
			user.Name,
			user.Email,
			user.Phone,
			user.Address,
			user.PostalCode,
			user.ProfileImageURL,
			user.Role.Name,
		)

		response.SuccessResponse(c, "profile retrieved successfull", profileResp)
	}
}

func UpdateProfile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint("user_id")
		var user models.User

		if err := db.First(&user, userID).Error; err != nil {
			response.ErrorResponse(c, http.StatusNotFound, "user not found")
			return
		}

		oldImagePublicID := user.ProfileImagePublicID
		contentType := c.GetHeader("Content-Type")
		updates := map[string]interface{}{}

		if strings.HasPrefix(contentType, "multipart/form-data") {
			if name := c.PostForm("name"); name != "" {
				updates["name"] = name
			}
			if phone := c.PostForm("phone"); phone != "" {
				updates["phone"] = phone
			}
			if address := c.PostForm("address"); address != "" {
				updates["address"] = address
			}
			if postal := c.PostForm("postal_code"); postal != "" {
				updates["postal_code"] = postal
			}

			file, err := c.FormFile("profile_image")
			if err == nil {
				tempPath := "/tmp/" + file.Filename
				if err := c.SaveUploadedFile(file, tempPath); err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed save image")
					return
				}

				upload, err := helper.UploadImageFromFile(tempPath, "user-profiles")
				os.Remove(tempPath)

				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "failed upload image")
					return
				}

				updates["profile_image_url"] = upload.SecureURL
				updates["profile_image_public_id"] = upload.PublicID
			}
		} else {
			var req request.UpdateProfileRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.ErrorResponse(c, http.StatusBadRequest, "invalid body")
				return
			}

			if req.Name != nil { updates["name"] = *req.Name }
			if req.Phone != nil { updates["phone"] = *req.Phone }
			if req.Address != nil { updates["address"] = *req.Address }
			if req.PostalCode != nil { updates["postal_code"] = *req.PostalCode }
		}

		if len(updates) == 0 {
			response.ErrorResponse(c, http.StatusBadRequest, "no data updated")
			return
		}

		db.Model(&user).Updates(updates)

		if oldImagePublicID != "" && updates["profile_image_public_id"] != nil {
			helper.DeleteImage(oldImagePublicID)
		}

		db.Preload("Role").First(&user, user.ID)

		resp := response.BuildUserProfileResponse(
			user.ID,
			user.Name,
			user.Email,
			user.Phone,
			user.Address,
			user.PostalCode,
			user.ProfileImageURL,
			user.Role.Name,
		)

		response.SuccessResponse(c, "profile updated", resp)
	}
}