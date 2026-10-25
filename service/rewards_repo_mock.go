package service

import (
	"context"
	entity "visa_assignment/models/entity"
)

type rewardRepoMock struct {
	balance int64
	rewards []entity.RewardEntity
	err     error
}

func (m *rewardRepoMock) Create(
	ctx context.Context,
	reward *entity.RewardEntity,
	isRedemption bool,
) error {
	if m.err != nil {
		return m.err
	}

	reward.RewardID = 1
	return nil
}

func (m *rewardRepoMock) GetRewardsByMemberID(
	ctx context.Context,
	memberID int64,
) ([]entity.RewardEntity, error) {
	return m.rewards, m.err
}

func (m *rewardRepoMock) GetMemberBalance(
	ctx context.Context,
	memberID int64,
) (int64, error) {
	return m.balance, m.err
}
