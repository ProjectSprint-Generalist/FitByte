package handlers

import (
	"net/http"
	"strconv"

	"fitbyte/internal/models"
	"fitbyte/internal/services"

	"github.com/gin-gonic/gin"
)


type UserHandler struct {
	userService *services.UserService
}


func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}


func (h *UserHandler) GetUsers(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))


	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userService.GetUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to retrieve users",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = h.userService.ToUserResponse(&user)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    userResponses,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() != "user not found" {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    statusCode,
		})
		return
	}


	userResponse := h.userService.ToUserResponse(user)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    userResponse,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	registerReq := &models.RegisterRequest{
		Email:      req.Email,
		Password:   "defaultPasswor",
		Name:       req.Name,
		Preference: req.Preference,
		WeightUnit: req.WeightUnit,
		HeightUnit: req.HeightUnit,
		Weight:     req.Weight,
		Height:     req.Height,
		ImageURI:   req.ImageURI,
	}

	user, err := h.userService.CreateUser(registerReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    statusCode,
		})
		return
	}

	userResponse := h.userService.ToUserResponse(user)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    userResponse,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() != "user not found" {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    statusCode,
		})
		return
	}

	userResponse := h.userService.ToUserResponse(user)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    userResponse,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to delete user",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
