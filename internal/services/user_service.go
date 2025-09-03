package services

import (
	"errors"
	"fmt"

	"fitbyte/internal/database"
	"fitbyte/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{
		db: database.DB,
	}
}


func (s *UserService) CreateUser(req *models.RegisterRequest) (*models.User, error) {

	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}


	user := &models.User{
		Email:      req.Email,
		Password:   string(hashedPassword),
		Name:       req.Name,
		Preference: req.Preference,
		WeightUnit: req.WeightUnit,
		HeightUnit: req.HeightUnit,
		Weight:     req.Weight,
		Height:     req.Height,
		ImageURI:   req.ImageURI,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}


func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}


func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}


func (s *UserService) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}


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

	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}


func (s *UserService) DeleteUser(id uint) error {
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}


func (s *UserService) GetUsers(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64


	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}


	offset := (page - 1) * limit


	if err := s.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}


func (s *UserService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}


func (s *UserService) ToUserResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Preference: user.Preference,
		WeightUnit: user.WeightUnit,
		HeightUnit: user.HeightUnit,
		Weight:     user.Weight,
		Height:     user.Height,
		ImageURI:   user.ImageURI,
	}
}
