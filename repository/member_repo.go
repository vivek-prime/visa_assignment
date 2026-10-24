package repository

import (
	"context"
	"database/sql"
	"errors"

	model "visa_assignment/models/entity"
	"visa_assignment/utils"
)

type MemberRepository interface {
	Create(ctx context.Context, member *model.MemberEntity) error
	GetByID(ctx context.Context, memberID int64) (*model.MemberEntity, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type memberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(ctx context.Context, member *model.MemberEntity) error {
	query := `
		INSERT INTO Members (name, email, created_at, updated_at)
		OUTPUT INSERTED.member_id
		VALUES (@p1, @p2, @p3, @p4)
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		member.Name,
		member.Email,
		member.CreatedAt,
		member.UpdatedAt,
	).Scan(&member.MemberID)

	if err != nil {
		return err
	}

	return nil
}

func (r *memberRepository) GetByID(ctx context.Context, memberID int64) (*model.MemberEntity, error) {
	query := `
		SELECT member_id, name, email, created_at, updated_at
		FROM Members
		WHERE Member_ID = @p1
	`

	var m model.MemberEntity
	err := r.db.QueryRowContext(ctx, query, memberID).Scan(
		&m.MemberID,
		&m.Name,
		&m.Email,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrMemberNotFound
		}
		return nil, err
	}

	return &m, nil
}

func (r *memberRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM Members
		WHERE Email = @p1
	`

	var count int
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}
