package usecase

import (
	"context"
	"fmt"

	"CalorieCompass/internal/entity"
	"CalorieCompass/internal/pkg/hash"
)

type UserRepo interface {
	Create(ctx context.Context, user entity.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
	GetByID(ctx context.Context, id int64) (entity.User, error)
}

type TokenRepo interface {
	GenerateToken(user entity.User) (string, error)
	ValidateToken(token string) (int64, error)
}

type AuthUseCase struct {
	userRepo  UserRepo
	tokenRepo TokenRepo
	hasher    *hash.Hasher
}

func NewAuthUseCase(userRepo UserRepo, tokenRepo TokenRepo, hasher *hash.Hasher) *AuthUseCase {
	return &AuthUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		hasher:    hasher,
	}
}

func (uc *AuthUseCase) SignUp(ctx context.Context, input entity.UserSignUp) (entity.AuthResponse, error) {
	_, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return entity.AuthResponse{}, fmt.Errorf("email already exists")
	}

	hashedPassword, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("password hash error: %w", err)
	}

	userID, err := uc.userRepo.Create(ctx, entity.User{
		Email:    input.Email,
		Password: hashedPassword,
		Name:     input.Name,
	})
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("create user error: %w", err)
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("get user error: %w", err)
	}

	token, err := uc.tokenRepo.GenerateToken(user)
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("generate token error: %w", err)
	}

	return entity.AuthResponse{
		User: entity.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		Token: token,
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input entity.UserLogin) (entity.AuthResponse, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("invalid email or password")
	}

	if !uc.hasher.Check(input.Password, user.Password) {
		return entity.AuthResponse{}, fmt.Errorf("invalid email or password")
	}

	token, err := uc.tokenRepo.GenerateToken(user)
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("generate token error: %w", err)
	}

	return entity.AuthResponse{
		User: entity.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		Token: token,
	}, nil
}

func (uc *AuthUseCase) GetUserByID(ctx context.Context, id int64) (entity.UserResponse, error) {
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