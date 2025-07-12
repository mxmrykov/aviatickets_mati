package usecase

import (
	"context"

	"aviatickets_mati/internal/entity"
	"aviatickets_mati/internal/repository"
)

type UserUseCase struct {
	userRepo *repository.UserRepository
}

func NewUserUseCase(userRepo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (uc *UserUseCase) GetOrCreate(ctx context.Context, telegramID int64, username, firstName, lastName string) (*entity.User, error) {
	user, err := uc.userRepo.FindByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	newUser := entity.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
	}

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return uc.userRepo.FindByTelegramID(ctx, telegramID)
}
