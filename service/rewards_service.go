package service

import (
	"context"
	"strings"
	"time"

	"rewards-ledger/model"
	"rewards-ledger/repository"
	"rewards-ledger/utils"
)

type RewardService interface {
	CreateReward(ctx context.Context, req model.CreateRewardRequest) (*model.Reward, error)
	GetRewardsByMemberID(ctx context.Context, memberID int64) ([]model.Reward, error)
}

type rewardService struct {
	memberRepo repository.MemberRepository
	rewardRepo repository.RewardRepository
}

func NewRewardService(memberRepo repository.MemberRepository, rewardRepo repository.RewardRepository) RewardService {
	return &rewardService{
		memberRepo: memberRepo,
		rewardRepo: rewardRepo,
	}
}

func (s *rewardService) CreateReward(ctx context.Context, req model.CreateRewardRequest) (*model.Reward, error) {
	// Rule: point_type_id must be a valid type (1-4)
	if !model.IsValidPointType(req.PointTypeID) {
		return nil, utils.ErrInvalidPointType
	}

	// Rule: points must be a positive number
	if req.Points <= 0 {
		return nil, utils.ErrInvalidPoints
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		return nil, utils.ErrInvalidInput
	}

	// Rule: member_id must reference an existing member
	_, err := s.memberRepo.GetByID(ctx, req.MemberID)
	if err != nil {
		return nil, err
	}

	// Rule: The system applies the sign based on point_type_id
	// Credits (1, 2, 3) are positive; Redemption (4) is negative
	signedPoints := req.Points
	if model.IsRedemption(req.PointTypeID) {
		signedPoints = -req.Points
	}

	reward := &model.Reward{
		MemberID:    req.MemberID,
		PointTypeID: req.PointTypeID,
		Points:      signedPoints,
		Description: description,
		EventDate:   time.Now().UTC(),
	}

	// For redemptions, CreateReward in the repository layer checks the member's balance
	// transactionally to ensure a redemption never exceeds available balance.
	if model.IsRedemption(req.PointTypeID) {
		if err := s.rewardRepo.CreateRedemption(ctx, reward, req.Points); err != nil {
			return nil, err
		}
		return reward, nil
	}

	if err := s.rewardRepo.Create(ctx, reward); err != nil {
		return nil, err
	}

	return reward, nil
}

func (s *rewardService) GetRewardsByMemberID(ctx context.Context, memberID int64) ([]model.Reward, error) {
	// Verify member exists first so invalid memberId returns 404 Not Found
	_, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	rewards, err := s.rewardRepo.GetByMemberID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	// Return an empty JSON array [] instead of null if the member has no reward entries yet
	if rewards == nil {
		return []model.Reward{}, nil
	}

	return rewards, nil
}
