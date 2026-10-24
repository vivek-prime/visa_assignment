package service

import (
	"context"
	"strings"
	"time"
	"visa_assignment/constants"
	dto "visa_assignment/models/dto"
	entity "visa_assignment/models/entity"
	"visa_assignment/repository"
	"visa_assignment/utils"
)

type RewardService interface {
	CreateReward(ctx context.Context, req dto.RewardRequest) (*entity.RewardEntity, error)
	GetRewardsByMemberID(ctx context.Context, memberID int64) ([]entity.RewardEntity, error)
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

func (s *rewardService) CreateReward(ctx context.Context, req dto.RewardRequest) (*entity.RewardEntity, error) {
	if err := validateRewardRequest(&req); err != nil {
		return nil, err
	}

	// member_id must reference an existing member
	_, err := s.memberRepo.GetByID(ctx, req.MemberID)
	if err != nil {
		return nil, err
	}

	//The system applies the sign based on point_type_id
	// Credits (1, 2, 3) are positive; Redemption (4) is negative
	signedPoints := req.Points
	if isRedemption(req.PointTypeID) {
		signedPoints = -req.Points
	}

	reward := &entity.RewardEntity{
		MemberID:    req.MemberID,
		PointTypeID: req.PointTypeID,
		Points:      signedPoints,
		Description: req.Description,
		EventDate:   time.Now().UTC(),
	}

	// For redemptions (point_type_id == 4), CreateRedemption locks the member row,
	// verifies member existence, checks available balance, and inserts atomically.
	if isRedemption(req.PointTypeID) {
		if err := s.rewardRepo.Create(ctx, reward, true); err != nil {
			return nil, err
		}
		return reward, nil
	}

	if err := s.rewardRepo.Create(ctx, reward, false); err != nil {
		return nil, err
	}

	return reward, nil
}

func (s *rewardService) GetRewardsByMemberID(ctx context.Context, memberID int64) ([]entity.RewardEntity, error) {
	// Verify member exists first so invalid memberId returns 404 Not Found
	_, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	rewards, err := s.rewardRepo.GetRewardsByMemberID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	// Return an empty JSON array [] instead of null if the member has no reward entries yet
	if rewards == nil {
		return []entity.RewardEntity{}, nil
	}

	return rewards, nil
}

func validateRewardRequest(req *dto.RewardRequest) error {
	if !isValidPointType(req.PointTypeID) {
		return utils.ErrInvalidPointType
	}

	if req.Points <= 0 {
		return utils.ErrInvalidPoints
	}

	req.Description = strings.TrimSpace(req.Description)
	if req.Description == constants.Empty {
		return utils.ErrInvalidInput
	}

	return nil
}
