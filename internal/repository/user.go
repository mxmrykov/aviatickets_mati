package repository

import (
	"context"
	_ "embed"
	"errors"

	"aviatickets_mati/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed sql/user/create.sql
	create string

	//go:embed sql/user/findByTgID.sql
	findByTgID string
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) error {
	args := pgx.NamedArgs{
		"telegram_id": user.TelegramID,
		"username":    user.Username,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
	}

	_, err := r.db.Exec(ctx, create, args)
	return err
}

func (r *UserRepository) FindByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error) {
	args := pgx.NamedArgs{
		"telegram_id": telegramID,
	}

	var user entity.User
	err := r.db.QueryRow(ctx, findByTgID, args).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
