package handlers

import (
	"net/http"
	"strconv"

	"fitbyte/internal/models"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	// In a real application, you would inject a service or repository here
	// userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUsers returns a list of users
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// In a real application, you would fetch users from a database
	name1 := "John Doe"
	name2 := "Jane Smith"
	preference1 := "metric"
	weightUnit1 := "kg"
	heightUnit1 := "cm"
	weight1 := 75.5
	height1 := 180.0
	imageURI1 := "https://example.com/image1.jpg"

	users := []models.UserResponse{
		{
			ID:         1,
			Email:      "john@example.com",
			Name:       &name1,
			Preference: &preference1,
			WeightUnit: &weightUnit1,
			HeightUnit: &heightUnit1,
			Weight:     &weight1,
			Height:     &height1,
			ImageURI:   &imageURI1,
		},
		{
			ID:    2,
			Email: "jane@example.com",
			Name:  &name2,
			// Other fields are nil (null in JSON)
		},
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int64(len(users)),
			TotalPages: 1,
		},
	})
}

// GetUser returns a specific user by ID
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

	// In a real application, you would fetch the user from a database
	name := "John Doe"
	preference := "metric"
	weightUnit := "kg"
	heightUnit := "cm"
	weight := 75.5
	height := 180.0
	imageURI := "https://example.com/image1.jpg"

	user := models.UserResponse{
		ID:         uint(id),
		Email:      "john@example.com",
		Name:       &name,
		Preference: &preference,
		WeightUnit: &weightUnit,
		HeightUnit: &heightUnit,
		Weight:     &weight,
		Height:     &height,
		ImageURI:   &imageURI,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// CreateUser creates a new user
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

	// In a real application, you would create the user in a database
	user := models.UserResponse{
		ID:         1,
		Email:      req.Email,
		Name:       req.Name,
		Preference: req.Preference,
		WeightUnit: req.WeightUnit,
		HeightUnit: req.HeightUnit,
		Weight:     req.Weight,
		Height:     req.Height,
		ImageURI:   req.ImageURI,
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser updates an existing user
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

	// In a real application, you would update the user in a database
	name := "John Doe"
	preference := "metric"
	weightUnit := "kg"
	heightUnit := "cm"
	weight := 75.5
	height := 180.0
	imageURI := "https://example.com/image1.jpg"

	user := models.UserResponse{
		ID:         uint(id),
		Email:      "john@example.com",
		Name:       &name,
		Preference: &preference,
		WeightUnit: &weightUnit,
		HeightUnit: &heightUnit,
		Weight:     &weight,
		Height:     &height,
		ImageURI:   &imageURI,
	}

	// Apply updates
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.Name = req.Name
	}
	if req.Preference != nil {
		user.Preference = req.Preference
	}
	if req.WeightUnit != nil {
		user.WeightUnit = req.WeightUnit
	}
	if req.HeightUnit != nil {
		user.HeightUnit = req.HeightUnit
	}
	if req.Weight != nil {
		user.Weight = req.Weight
	}
	if req.Height != nil {
		user.Height = req.Height
	}
	if req.ImageURI != nil {
		user.ImageURI = req.ImageURI
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser deletes a user
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

	// In a real application, you would delete the user from a database
	_ = id // Use the ID to delete the user

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
