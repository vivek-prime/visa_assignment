package service

import (
	"context"
	"errors"
	"testing"

	dto "visa_assignment/models/dto"
	entity "visa_assignment/models/entity"
	"visa_assignment/utils"
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

func TestCreateReward_CreditSuccess(t *testing.T) {
	memberRepo := &memberRepoMock{
		member: &entity.MemberEntity{
			MemberID: 1,
			Name:     "Alex",
		},
	}

	rewardRepo := &rewardRepoMock{}

	svc := NewRewardService(memberRepo, rewardRepo)

	reward, err := svc.CreateReward(context.Background(), dto.RewardRequest{
		MemberID:    1,
		PointTypeID: 1,
		Points:      100,
		Description: "  Signup bonus  ",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if reward.Points != 100 {
		t.Errorf("expected 100 points, got %d", reward.Points)
	}

	if reward.Description != "Signup bonus" {
		t.Errorf("expected trimmed description, got %q", reward.Description)
	}

	if reward.MemberID != 1 {
		t.Errorf("expected member ID 1, got %d", reward.MemberID)
	}

	if reward.EventDate.IsZero() {
		t.Error("expected EventDate to be set")
	}
}

func TestCreateReward_RedemptionSuccess(t *testing.T) {
	memberRepo := &memberRepoMock{
		member: &entity.MemberEntity{MemberID: 1},
	}

	rewardRepo := &rewardRepoMock{}

	svc := NewRewardService(memberRepo, rewardRepo)

	reward, err := svc.CreateReward(context.Background(), dto.RewardRequest{
		MemberID:    1,
		PointTypeID: 4,
		Points:      50,
		Description: "Redeemed reward",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if reward.Points != -50 {
		t.Errorf("expected -50 points for redemption, got %d", reward.Points)
	}
}

func TestCreateReward_InvalidPointType(t *testing.T) {
	svc := NewRewardService(&memberRepoMock{}, &rewardRepoMock{})

	reward, err := svc.CreateReward(context.Background(), dto.RewardRequest{
		MemberID:    1,
		PointTypeID: 99,
		Points:      100,
		Description: "Test reward",
	})

	if reward != nil {
		t.Error("expected nil reward")
	}

	if !errors.Is(err, utils.ErrInvalidPointType) {
		t.Errorf("expected ErrInvalidPointType, got %v", err)
	}
}

func TestCreateReward_InvalidPoints(t *testing.T) {
	svc := NewRewardService(&memberRepoMock{}, &rewardRepoMock{})

	reward, err := svc.CreateReward(context.Background(), dto.RewardRequest{
		MemberID:    1,
		PointTypeID: 1,
		Points:      0,
		Description: "Test reward",
	})

	if reward != nil {
		t.Error("expected nil reward")
	}

	if !errors.Is(err, utils.ErrInvalidPoints) {
		t.Errorf("expected ErrInvalidPoints, got %v", err)
	}
}

func TestCreateReward_EmptyDescription(t *testing.T) {
	svc := NewRewardService(&memberRepoMock{}, &rewardRepoMock{})

	reward, err := svc.CreateReward(context.Background(), dto.RewardRequest{
		MemberID:    1,
		PointTypeID: 1,
		Points:      100,
		Description: "   ",
	})

	if reward != nil {
		t.Error("expected nil reward")
	}

	if !errors.Is(err, utils.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreateReward_MemberNotFound(t *testing.T) {
	rewardRepo := &rewardRepoMock{
		err: utils.ErrMemberNotFound,
	}
	svc := NewRewardService(&memberRepoMock{}, rewardRepo)

	reward, err := svc.CreateReward(context.Background(), dto.RewardRequest{
		MemberID:    99,
		PointTypeID: 4,
		Points:      100,
		Description: "Test reward",
	})

	if reward != nil {
		t.Error("expected nil reward")
	}

	if !errors.Is(err, utils.ErrMemberNotFound) {
		t.Errorf("expected ErrMemberNotFound, got %v", err)
	}
}

func TestCreateReward_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	memberRepo := &memberRepoMock{
		member: &entity.MemberEntity{MemberID: 1},
	}

	rewardRepo := &rewardRepoMock{
		err: expectedErr,
	}

	svc := NewRewardService(memberRepo, rewardRepo)

	reward, err := svc.CreateReward(context.Background(), dto.RewardRequest{
		MemberID:    1,
		PointTypeID: 1,
		Points:      100,
		Description: "Test reward",
	})

	if reward != nil {
		t.Error("expected nil reward")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected database error, got %v", err)
	}
}

func TestGetRewardsByMemberID_Success(t *testing.T) {
	memberRepo := &memberRepoMock{
		member: &entity.MemberEntity{MemberID: 1},
	}

	rewardRepo := &rewardRepoMock{
		rewards: []entity.RewardEntity{
			{
				MemberID:    1,
				PointTypeID: 1,
				Points:      100,
				Description: "Signup bonus",
			},
			{
				MemberID:    1,
				PointTypeID: 4,
				Points:      -20,
				Description: "Redemption",
			},
		},
	}

	svc := NewRewardService(memberRepo, rewardRepo)

	rewards, err := svc.GetRewardsByMemberID(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(rewards) != 2 {
		t.Errorf("expected 2 rewards, got %d", len(rewards))
	}

	if rewards[0].Points != 100 {
		t.Errorf("expected first reward to have 100 points, got %d", rewards[0].Points)
	}

	if rewards[1].Points != -20 {
		t.Errorf("expected second reward to have -20 points, got %d", rewards[1].Points)
	}
}

func TestGetRewardsByMemberID_EmptyList(t *testing.T) {
	memberRepo := &memberRepoMock{
		member: &entity.MemberEntity{MemberID: 1},
	}

	rewardRepo := &rewardRepoMock{}

	svc := NewRewardService(memberRepo, rewardRepo)

	rewards, err := svc.GetRewardsByMemberID(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if rewards == nil {
		t.Error("expected empty slice, got nil")
	}

	if len(rewards) != 0 {
		t.Errorf("expected 0 rewards, got %d", len(rewards))
	}
}

func TestGetRewardsByMemberID_MemberNotFound(t *testing.T) {
	memberRepo := &memberRepoMock{
		err: utils.ErrMemberNotFound,
	}

	svc := NewRewardService(memberRepo, &rewardRepoMock{})

	rewards, err := svc.GetRewardsByMemberID(context.Background(), 99)

	if rewards != nil {
		t.Error("expected nil rewards")
	}

	if !errors.Is(err, utils.ErrMemberNotFound) {
		t.Errorf("expected ErrMemberNotFound, got %v", err)
	}
}
