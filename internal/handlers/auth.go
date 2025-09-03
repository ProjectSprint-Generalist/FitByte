package handlers

import (
	"net/http"

	"fitbyte/internal/models"
	"fitbyte/internal/services"

	"github.com/gin-gonic/gin"
)


type AuthHandler struct {
	userService *services.UserService
	jwtService  *services.JWTService
}


func NewAuthHandler(userService *services.UserService, jwtService *services.JWTService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
	}
}


func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}


	user, err := h.userService.CreateUser(&req)
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


	token, err := h.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to generate authentication token",
			Code:    http.StatusInternalServerError,
		})
		return
	}


	userResponse := h.userService.ToUserResponse(user)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data: models.AuthResponse{
			User:  userResponse,
			Token: token,
		},
	})
}


func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}


	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Invalid email or password",
			Code:    http.StatusUnauthorized,
		})
		return
	}


	if err := h.userService.VerifyPassword(user.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Invalid email or password",
			Code:    http.StatusUnauthorized,
		})
		return
	}


	token, err := h.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to generate authentication token",
			Code:    http.StatusInternalServerError,
		})
		return
	}


	userResponse := h.userService.ToUserResponse(user)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: models.AuthResponse{
			User:  userResponse,
			Token: token,
		},
	})
}


func (h *AuthHandler) Profile(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}


	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}


	userResponse := h.userService.ToUserResponse(user)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    userResponse,
	})
}


func (h *AuthHandler) RefreshToken(c *gin.Context) {

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Authorization header is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	tokenString := authHeader[7:]


	newToken, err := h.jwtService.RefreshToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Invalid or expired token",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: gin.H{
			"token": newToken,
		},
	})
}
