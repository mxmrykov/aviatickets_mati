package usecase

import (
	"context"

	"aviatickets_mati/internal/entity"
	"aviatickets_mati/internal/repository"
)

type BookingUseCase struct {
	bookingRepo *repository.BookingRepository
	flightRepo  *repository.FlightRepository
}

func NewBookingUseCase(
	bookingRepo *repository.BookingRepository,
	flightRepo *repository.FlightRepository,
) *BookingUseCase {
	return &BookingUseCase{
		bookingRepo: bookingRepo,
		flightRepo:  flightRepo,
	}
}

func (uc *BookingUseCase) CreateBooking(ctx context.Context, booking entity.Booking) error {
	return uc.bookingRepo.Create(ctx, booking)
}

func (uc *BookingUseCase) GetUserBookings(ctx context.Context, userID int) ([]entity.Booking, error) {
	return uc.bookingRepo.FindByUserID(ctx, userID)
}

func (uc *BookingUseCase) GetBookingByID(ctx context.Context, id int) (*entity.Booking, error) {
	return uc.bookingRepo.FindByID(ctx, id)
}

func (uc *BookingUseCase) CancelBooking(ctx context.Context, id int) error {
	return uc.bookingRepo.Cancel(ctx, id)
}
