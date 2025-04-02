package usecase

import (
	"context"
	"fmt"

	"CalorieCompass/internal/entity"
)

// UserUseCase - kullanıcı use case'i
type UserUseCase struct {
	userRepo UserRepo
}

// NewUserUseCase - yeni bir user use case oluşturur
func NewUserUseCase(userRepo UserRepo) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// GetUserByID - ID'ye göre kullanıcı bilgisini getirir
func (uc *UserUseCase) GetUserByID(ctx context.Context, id int64) (entity.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return entity.UserResponse{}, fmt.Errorf("get user error: %w", err)
	}

	return entity.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}