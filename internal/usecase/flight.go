package usecase

import (
	"context"

	"aviatickets_mati/internal/entity"
	"aviatickets_mati/internal/repository"
)

type FlightUseCase struct {
	flightRepo *repository.FlightRepository
}

func NewFlightUseCase(flightRepo *repository.FlightRepository) *FlightUseCase {
	return &FlightUseCase{flightRepo: flightRepo}
}

func (uc *FlightUseCase) GetAllFlights(ctx context.Context) ([]entity.Flight, error) {
	return uc.flightRepo.FindAll(ctx)
}

func (uc *FlightUseCase) GetFlightByID(ctx context.Context, id int) (*entity.Flight, error) {
	return uc.flightRepo.FindById(ctx, id)
}

func (uc *FlightUseCase) SearchFlights(ctx context.Context, origin, destination string) ([]entity.Flight, error) {
	return uc.flightRepo.Search(ctx, origin, destination)
}

func (uc *FlightUseCase) CreateFlight(ctx context.Context, flight entity.Flight) error {
	return uc.flightRepo.Create(ctx, flight)
}

func (uc *FlightUseCase) DeleteFlight(ctx context.Context, id int) error {
	return uc.flightRepo.Delete(ctx, id)
}

func (uc *FlightUseCase) GetFlightStats(ctx context.Context, flightID int) (*entity.FlightStats, error) {
	return uc.flightRepo.GetStats(ctx, flightID)
}

func (uc *FlightUseCase) GetFlightByNumber(ctx context.Context, flightNumber string) (*entity.Flight, error) {
	return uc.flightRepo.FindByFlightNumber(ctx, flightNumber)
}
