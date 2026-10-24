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
}

type memberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(ctx context.Context, member *model.MemberEntity) error {
	query := `
		INSERT INTO vivek_db.dbo.members (name, email, created_at)
		OUTPUT INSERTED.member_id
		VALUES (@p1, @p2, @p3)
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		member.Name,
		member.Email,
		member.CreatedAt,
	).Scan(&member.MemberID)

	if err != nil {
		return err
	}

	return nil
}

func (r *memberRepository) GetByID(ctx context.Context, memberID int64) (*model.MemberEntity, error) {
	query := `
		SELECT member_id, name, email, created_at
		FROM vivek_db.dbo.members
		WHERE Member_ID = @p1
	`

	var m model.MemberEntity
	err := r.db.QueryRowContext(ctx, query, memberID).Scan(
		&m.MemberID,
		&m.Name,
		&m.Email,
		&m.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrMemberNotFound
		}
		return nil, err
	}

	return &m, nil
}
