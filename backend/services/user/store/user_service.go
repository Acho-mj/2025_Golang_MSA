package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"Acho-mj/2025_Golang_MSA/backend/internal/storage"
	"Acho-mj/2025_Golang_MSA/backend/services/user/models"
)

var (
	ErrInvalidInput = errors.New("잘못된 입력입니다")
	ErrUserNotFound = errors.New("사용자를 찾을 수 없습니다")
)

type UserService struct {
	storage *storage.UserStorage
}

func NewUserService(storage *storage.UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*models.User, error) {
	if email == "" || name == "" {
		return nil, fmt.Errorf("%w: email과 name은 필수입니다", ErrInvalidInput)
	}

	item := &storage.UserItem{
		UserID:    generateUserID(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.storage.CreateUser(ctx, item); err != nil {
		return nil, err
	}

	return &models.User{
		UserID:    item.UserID,
		Email:     item.Email,
		Name:      item.Name,
		CreatedAt: item.CreatedAt,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("%w: userID는 필수입니다", ErrInvalidInput)
	}

	item, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.User{
		UserID:    item.UserID,
		Email:     item.Email,
		Name:      item.Name,
		CreatedAt: item.CreatedAt,
	}, nil
}

func generateUserID() string {
	return fmt.Sprintf("user-%d", time.Now().UnixNano())
}
