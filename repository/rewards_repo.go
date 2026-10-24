package repository

import (
	"context"
	"database/sql"
	"errors"

	model "visa_assignment/models/entity"
	"visa_assignment/utils"
)

type RewardRepository interface {
	Create(ctx context.Context, reward *model.RewardEntity) error
	CreateRedemption(ctx context.Context, reward *model.RewardEntity, redemptionAmount int64) error
	GetByMemberID(ctx context.Context, memberID int64) ([]model.RewardEntity, error)
	GetMemberBalance(ctx context.Context, memberID int64) (int64, error)
}

type rewardRepository struct {
	db *sql.DB
}

func NewRewardRepository(db *sql.DB) RewardRepository {
	return &rewardRepository{db: db}
}

func (r *rewardRepository) Create(ctx context.Context, reward *model.RewardEntity) error {
	query := `
		INSERT INTO Rewards (Member_ID, PointType_ID, Points, Description, EventDate)
		OUTPUT INSERTED.Reward_ID
		VALUES (@p1, @p2, @p3, @p4, @p5)
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		reward.MemberID,
		reward.PointTypeID,
		reward.Points,
		reward.Description,
		reward.EventDate,
	).Scan(&reward.RewardID)
}

// CreateRedemption checks available balance and inserts the negative reward entry atomically
func (r *rewardRepository) CreateRedemption(ctx context.Context, reward *model.RewardEntity, redemptionAmount int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Lock the member row for this transaction so concurrent redemptions serialize cleanly
	lockMemberQuery := `
		SELECT Member_ID
		FROM Members WITH (UPDLOCK, HOLDLOCK)
		WHERE Member_ID = @p1
	`
	var lockedMemberID int64
	if err := tx.QueryRowContext(ctx, lockMemberQuery, reward.MemberID).Scan(&lockedMemberID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrMemberNotFound
		}
		return err
	}

	// Calculate current available balance inside the transaction
	balanceQuery := `
		SELECT COALESCE(SUM(Points), 0)
		FROM Rewards
		WHERE Member_ID = @p1
	`
	var currentBalance int64
	if err := tx.QueryRowContext(ctx, balanceQuery, reward.MemberID).Scan(&currentBalance); err != nil {
		return err
	}

	// Rule: A redemption must not exceed the member's available balance
	if currentBalance < redemptionAmount {
		return utils.ErrInsufficientBalance
	}

	insertQuery := `
		INSERT INTO Rewards (Member_ID, PointType_ID, Points, Description, EventDate)
		OUTPUT INSERTED.Reward_ID
		VALUES (@p1, @p2, @p3, @p4, @p5)
	`
	err = tx.QueryRowContext(
		ctx,
		insertQuery,
		reward.MemberID,
		reward.PointTypeID,
		reward.Points,
		reward.Description,
		reward.EventDate,
	).Scan(&reward.RewardID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *rewardRepository) GetByMemberID(ctx context.Context, memberID int64) ([]model.RewardEntity, error) {
	query := `
		SELECT Reward_ID, Member_ID, PointType_ID, Points, Description, EventDate
		FROM Rewards
		WHERE Member_ID = @p1
		ORDER BY EventDate ASC, Reward_ID ASC
	`

	rows, err := r.db.QueryContext(ctx, query, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rewards []model.RewardEntity
	for rows.Next() {
		var rw model.RewardEntity
		if err := rows.Scan(
			&rw.RewardID,
			&rw.MemberID,
			&rw.PointTypeID,
			&rw.Points,
			&rw.Description,
			&rw.EventDate,
		); err != nil {
			return nil, err
		}
		rewards = append(rewards, rw)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rewards, nil
}

func (r *rewardRepository) GetMemberBalance(ctx context.Context, memberID int64) (int64, error) {
	query := `
		SELECT COALESCE(SUM(Points), 0)
		FROM Rewards
		WHERE Member_ID = @p1
	`

	var balance int64
	if err := r.db.QueryRowContext(ctx, query, memberID).Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}
