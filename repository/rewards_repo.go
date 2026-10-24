package repository

import (
	"context"
	"database/sql"
	"errors"
	utils "visa_assignment/utils"

	model "visa_assignment/models/entity"

	mssql "github.com/microsoft/go-mssqldb"
)

type RewardRepository interface {
	Create(ctx context.Context, reward *model.RewardEntity, isRedemption bool) error
	GetRewardsByMemberID(ctx context.Context, memberID int64) ([]model.RewardEntity, error)
	GetMemberBalance(ctx context.Context, memberID int64) (int64, error)
}

type rewardRepository struct {
	db *sql.DB
}

func NewRewardRepository(db *sql.DB) RewardRepository {
	return &rewardRepository{db: db}
}

func (r *rewardRepository) Create(ctx context.Context, reward *model.RewardEntity, isRedemption bool) error {
	if !isRedemption {
		return r.createReward(ctx, reward)
	}
	return r.createRedemption(ctx, reward)
}

func (r *rewardRepository) createReward(ctx context.Context, reward *model.RewardEntity) error {
	query := `
		INSERT INTO vivek_db.dbo.rewards (Member_ID, PointType_ID, Points, Description, EventDate)
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

func (r *rewardRepository) createRedemption(ctx context.Context, reward *model.RewardEntity) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// ROWLOCK: Lock only this specific member row
	// UPDLOCK: Hold update lock until transaction completes
	// NOWAIT: Fail immediately with SQL Error 1222 if another redemption is already locking this member
	lockMemberQuery := `
		SELECT member_id
		FROM vivek_db.dbo.members WITH (ROWLOCK, UPDLOCK, NOWAIT)
		WHERE member_id = @p1
	`
	var lockedMemberID int64
	if err := tx.QueryRowContext(ctx, lockMemberQuery, reward.MemberID).Scan(&lockedMemberID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrMemberNotFound
		}
		// SQL Server Error 1222 = Lock request time out period exceeded (triggered by NOWAIT)
		var sqlErr mssql.Error
		if errors.As(err, &sqlErr) && sqlErr.Number == 1222 {
			return utils.ErrConcurrentRedemption
		}
		return err
	}

	// Calculate current available balance inside the transaction
	balanceQuery := `
		SELECT COALESCE(SUM(points), 0)
		FROM vivek_db.dbo.rewards
		WHERE member_id = @p1
	`
	var currentBalance int64
	if err := tx.QueryRowContext(ctx, balanceQuery, reward.MemberID).Scan(&currentBalance); err != nil {
		return err
	}

	redemptionAmount := -reward.Points
	if currentBalance < redemptionAmount {
		return utils.ErrInsufficientBalance
	}

	insertQuery := `
		INSERT INTO rewards (member_id, point_type_id, points, description, event_date)
		OUTPUT INSERTED.reward_id
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

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *rewardRepository) GetRewardsByMemberID(ctx context.Context, memberID int64) ([]model.RewardEntity, error) {
	query := `
		SELECT Reward_ID, Member_ID, PointType_ID, Points, Description, EventDate
		FROM vivek_db.dbo.rewards
		WHERE Member_ID = @p1
		ORDER BY EventDate ASC, Reward_ID ASC
	`

	rows, err := r.db.QueryContext(ctx, query, memberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
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
		FROM vivek_db.dbo.rewards
		WHERE Member_ID = @p1
	`

	var balance int64
	if err := r.db.QueryRowContext(ctx, query, memberID).Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}
