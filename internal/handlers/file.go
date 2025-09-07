package handlers

import (
	"net/http"
	"strconv"

	"fitbyte/internal/models"
	"fitbyte/internal/services"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	minioService *services.MinIOService
}

func NewFileHandler(minioService *services.MinIOService) *FileHandler {
	return &FileHandler{
		minioService: minioService,
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "File is required",
		})
		return
	}

	response, err := h.minioService.UploadFile(file, userIDUint)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *FileHandler) GetUserFiles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	files, err := h.minioService.GetUserFiles(userIDUint, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve files",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "Files retrieved successfully",
		Data:    files,
	})
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	fileURI := c.Query("uri")
	if fileURI == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "File URI is required",
		})
		return
	}

	err := h.minioService.DeleteFile(fileURI, userIDUint)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "File deleted successfully",
	})
}
