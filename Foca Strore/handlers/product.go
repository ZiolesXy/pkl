package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"voca-store/helper"
	"voca-store/models"
	"voca-store/request"
	"voca-store/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if it's multipart form using strings.HasPrefix
		contentType := c.GetHeader("Content-Type")
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")
		
		if isMultipart {
			// Handle multipart form data
			name := c.PostForm("name")
			description := c.PostForm("description")
			priceStr := c.PostForm("price")
			stockStr := c.PostForm("stock")
			
			if name == "" {
				response.ErrorResponse(c, http.StatusBadRequest, "Name is required")
				return
			}
			
			// Parse price
			var price float64
			if priceStr != "" {
				_, err := fmt.Sscanf(priceStr, "%f", &price)
				if err != nil {
					response.ErrorResponse(c, http.StatusBadRequest, "Invalid price format")
					return
				}
			}
			
			// Parse stock
			var stock int
			if stockStr != "" {
				_, err := fmt.Sscanf(stockStr, "%d", &stock)
				if err != nil {
					response.ErrorResponse(c, http.StatusBadRequest, "Invalid stock format")
					return
				}
			}
			
			// Handle file upload
			var imageURL, imagePublicID string
			file, err := c.FormFile("image")
			if err == nil && file != nil {
				// Save file temporarily
				tempPath := "/tmp/" + file.Filename
				if err := c.SaveUploadedFile(file, tempPath); err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to save uploaded file")
					return
				}
				
				// Upload to Cloudinary
				uploadResult, err := helper.UploadImageFromFile(tempPath, "products")
				if err != nil {
					os.Remove(tempPath)
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload image")
					return
				}
				
				imageURL = uploadResult.SecureURL
				imagePublicID = uploadResult.PublicID // ✅ SIMPAN PUBLIC ID
				
				// Clean up temp file
				os.Remove(tempPath)
			} else if c.PostForm("image_url") != "" {
				// Upload from URL
				uploadResult, err := helper.UploadImageFromURL(c.PostForm("image_url"), "products")
				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload image from URL")
					return
				}
				imageURL = uploadResult.SecureURL
				imagePublicID = uploadResult.PublicID // ✅ SIMPAN PUBLIC ID
			}

			// Create product with ImagePublicID
			product := models.Product{
				Name:        name,
				Description: description,
				ImageURL:    imageURL,
				ImagePublicID: imagePublicID, // ✅ INI YANG DIPERBAIKI
				Price:       price,
				Stock:       stock,
			}

			if err := db.Create(&product).Error; err != nil {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product")
				return
			}

			// Build product response (without ImagePublicID)
			productResp := response.BuildProductResponse(
				product.ID,
				product.Name,
				product.Description,
				product.ImageURL,
				product.Price,
				product.Stock,
				product.CreatedAt,
				product.UpdatedAt,
			)

			response.SuccessResponse(c, "Product created successfully", productResp)
		} else {
			// Handle JSON data
			var req request.CreateProductRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
				return
			}
			
			// If imageURL is provided and is a URL, upload to Cloudinary
			var imageURL, imagePublicID string
			if req.ImageURL != "" {
				if len(req.ImageURL) > 4 && (req.ImageURL[:4] == "http" || req.ImageURL[:5] == "https") {
					uploadResult, err := helper.UploadImageFromURL(req.ImageURL, "products")
					if err != nil {
						response.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload image from URL")
						return
					}
					imageURL = uploadResult.SecureURL
					imagePublicID = uploadResult.PublicID // ✅ SIMPAN PUBLIC ID
				} else {
					// Not a URL, treat as plain string (no upload)
					imageURL = req.ImageURL
					imagePublicID = "" // No public ID for non-Cloudinary URLs
				}
			}

			// Create product with ImagePublicID
			product := models.Product{
				Name:        req.Name,
				Description: req.Description,
				ImageURL:    imageURL,
				ImagePublicID: imagePublicID, // ✅ INI YANG DIPERBAIKI
				Price:       req.Price,
				Stock:       req.Stock,
			}

			if err := db.Create(&product).Error; err != nil {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to create product")
				return
			}

			// Build product response
			productResp := response.BuildProductResponse(
				product.ID,
				product.Name,
				product.Description,
				product.ImageURL,
				product.Price,
				product.Stock,
				product.CreatedAt,
				product.UpdatedAt,
			)

			response.SuccessResponse(c, "Product created successfully", productResp)
		}
	}
}

func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		var product models.Product
		if err := db.First(&product, productID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(c, http.StatusNotFound, "Product not found")
			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch product")
			}
			return
		}

		// Store old image info for cleanup
		// oldImageURL := product.ImageURL
		oldImagePublicID := product.ImagePublicID

		// Handle multipart form for image upload
		contentType := c.GetHeader("Content-Type")
		isMultipart := strings.HasPrefix(contentType, "multipart/form-data")
		
		var newImagePublicID string
		updates := make(map[string]interface{})
		
		if isMultipart {
			// Handle multipart form data
			if name := c.PostForm("name"); name != "" {
				updates["name"] = name
			}
			if description := c.PostForm("description"); description != "" {
				updates["description"] = description
			}
			if priceStr := c.PostForm("price"); priceStr != "" {
				var price float64
				_, err := fmt.Sscanf(priceStr, "%f", &price)
				if err != nil {
					response.ErrorResponse(c, http.StatusBadRequest, "Invalid price format")
					return
				}
				updates["price"] = price
			}
			if stockStr := c.PostForm("stock"); stockStr != "" {
				var stock int
				_, err := fmt.Sscanf(stockStr, "%d", &stock)
				if err != nil {
					response.ErrorResponse(c, http.StatusBadRequest, "Invalid stock format")
					return
				}
				updates["stock"] = stock
			}
			
			// Handle file upload
			file, err := c.FormFile("image")
			if err == nil && file != nil {
				// Save file temporarily
				tempPath := "/tmp/" + file.Filename
				if err := c.SaveUploadedFile(file, tempPath); err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to save uploaded file")
					return
				}
				
				// Upload to Cloudinary
				uploadResult, err := helper.UploadImageFromFile(tempPath, "products")
				if err != nil {
					os.Remove(tempPath)
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload image")
					return
				}
				
				updates["image_url"] = uploadResult.SecureURL
				updates["image_public_id"] = uploadResult.PublicID
				newImagePublicID = uploadResult.PublicID
				
				// Clean up temp file
				os.Remove(tempPath)
			} else if c.PostForm("image_url") != "" {
				// Upload from URL
				uploadResult, err := helper.UploadImageFromURL(c.PostForm("image_url"), "products")
				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload image from URL")
					return
				}
				updates["image_url"] = uploadResult.SecureURL
				updates["image_public_id"] = uploadResult.PublicID
				newImagePublicID = uploadResult.PublicID
			}
		} else {
			// Handle JSON data
			var req request.UpdateProductRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				response.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
				return
			}

			if req.Name != nil {
				updates["name"] = *req.Name
			}
			if req.Description != nil {
				updates["description"] = *req.Description
			}
			if req.Price != nil {
				updates["price"] = *req.Price
			}
			if req.Stock != nil {
				updates["stock"] = *req.Stock
			}
			
			// Handle image URL update
			if req.ImageURL != nil && *req.ImageURL != "" {
				// Upload from URL
				uploadResult, err := helper.UploadImageFromURL(*req.ImageURL, "products")
				if err != nil {
					response.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload image from URL")
					return
				}
				updates["image_url"] = uploadResult.SecureURL
				updates["image_public_id"] = uploadResult.PublicID
				newImagePublicID = uploadResult.PublicID
			}
		}

		if len(updates) == 0 {
			response.ErrorResponse(c, http.StatusBadRequest, "No fields to update")
			return
		}

		if err := db.Model(&product).Updates(updates).Error; err != nil {
			// If update failed but new image was uploaded, clean up Cloudinary
			if newImagePublicID != "" {
				helper.DeleteImage(newImagePublicID)
			}
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to update product")
			return
		}

		// DELETE OLD IMAGE FROM CLOUDINARY IF IT WAS UPDATED
		if oldImagePublicID != "" && newImagePublicID != "" {
			// Delete the old image since we have a new one
			if err := helper.DeleteImage(oldImagePublicID); err != nil {
				// Log error but don't fail the entire operation
				fmt.Printf("Warning: Failed to delete old image %s: %v\n", oldImagePublicID, err)
			}
		} else if oldImagePublicID != "" {
			// Check if image was explicitly cleared
			if isMultipart {
				if c.PostForm("image_url") == "" {
					if err := helper.DeleteImage(oldImagePublicID); err != nil {
						fmt.Printf("Warning: Failed to delete old image %s: %v\n", oldImagePublicID, err)
					}
				}
			} else {
				// For JSON, check if req exists and image_url is empty
				var reqForCheck request.UpdateProductRequest
				if c.ShouldBindJSON(&reqForCheck) == nil && reqForCheck.ImageURL != nil && *reqForCheck.ImageURL == "" {
					if err := helper.DeleteImage(oldImagePublicID); err != nil {
						fmt.Printf("Warning: Failed to delete old image %s: %v\n", oldImagePublicID, err)
					}
				}
			}
		}

		// Reload product
		if err := db.First(&product, product.ID).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to reload product")
			return
		}

		// Build product response
		productResp := response.BuildProductResponse(
			product.ID,
			product.Name,
			product.Description,
			product.ImageURL,
			product.Price,
			product.Stock,
			product.CreatedAt,
			product.UpdatedAt,
		)

		response.SuccessResponse(c, "Product updated successfully", productResp)
	}
}

func DeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")
		var product models.Product
		if err := db.First(&product, productID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.ErrorResponse(c, http.StatusNotFound, "Product not found")
			} else {
				response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch product")
			}
			return
		}

		// Delete image from Cloudinary if exists
		if product.ImagePublicID != "" {
			if err := helper.DeleteImage(product.ImagePublicID); err != nil {
				fmt.Printf("Warning: Failed to delete image %s: %v\n", product.ImagePublicID, err)
			}
		}

		if err := db.Delete(&product).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete product")
			return
		}

		response.SuccessResponse(c, "Product deleted successfully", nil)
	}
}

func GetProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product
		if err := db.Find(&products).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch products")
			return
		}

		// Build product responses
		var productResponses []response.ProductResponse
		for _, p := range products {
			productResp := response.BuildProductResponse(
				p.ID,
				p.Name,
				p.Description,
				p.ImageURL,
				p.Price,
				p.Stock,
				p.CreatedAt,
				p.UpdatedAt,
			)
			productResponses = append(productResponses, productResp)
		}

		productListResp := response.BuildProductListResponse(productResponses)
		response.SuccessListResponse(c, "Products retrieved successfully", productListResp)
	}
}

func GetAllProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product
		if err := db.Find(&products).Error; err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch products")
			return
		}

		// Build product responses
		var productResponses []response.ProductResponse
		for _, p := range products {
			productResp := response.BuildProductResponse(
				p.ID,
				p.Name,
				p.Description,
				p.ImageURL,
				p.Price,
				p.Stock,
				p.CreatedAt,
				p.UpdatedAt,
			)
			productResponses = append(productResponses, productResp)
		}
		response.SuccessListResponse(c, "Products retrieved successfully", productResponses)
	}
}