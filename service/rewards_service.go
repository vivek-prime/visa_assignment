package service

import (
	"context"
	"strings"

	"visa_assignment/config"
	"visa_assignment/constants"
	dto "visa_assignment/models/dto"
	entity "visa_assignment/models/entity"
	"visa_assignment/repository"
	"visa_assignment/utils"

	"github.com/rs/zerolog"
)

type RewardService interface {
	CreateReward(ctx context.Context, req dto.RewardRequest) (*entity.RewardEntity, error)
	GetRewardsByMemberID(ctx context.Context, memberID int64) ([]entity.RewardEntity, error)
}

type rewardService struct {
	memberRepo repository.MemberRepository
	rewardRepo repository.RewardRepository
	pointTypes map[int]config.PointTypeConfig
}

func NewRewardService(
	memberRepo repository.MemberRepository,
	rewardRepo repository.RewardRepository,
	pointTypes map[int]config.PointTypeConfig,
) RewardService {
	return &rewardService{
		memberRepo: memberRepo,
		rewardRepo: rewardRepo,
		pointTypes: pointTypes,
	}
}

func (s *rewardService) CreateReward(ctx context.Context, req dto.RewardRequest) (*entity.RewardEntity, error) {
	if err := s.validateRewardRequest(ctx, &req); err != nil {
		return nil, err
	}

	// Credit entries are positive; debit entries are negative.
	signedPoints := req.Points
	if isPointTypeDebit(req.PointTypeID, s.pointTypes) {
		signedPoints = -req.Points
	}

	reward := &entity.RewardEntity{
		MemberID:    req.MemberID,
		PointTypeID: req.PointTypeID,
		Points:      signedPoints,
		Description: req.Description,
		EventDate:   utils.GetCurrentTime(),
	}

	// Debit entries lock the member row, verify member existence and balance, and insert atomically.
	if s.pointTypes[req.PointTypeID] == config.PointTypeDebit {
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

func (s *rewardService) validateRewardRequest(ctx context.Context, req *dto.RewardRequest) error {
	if !isValidPointType(req.PointTypeID, s.pointTypes) {
		return utils.ErrInvalidPointType
	}

	if req.Points <= 0 {
		return utils.ErrInvalidPoints
	}

	req.Description = strings.TrimSpace(req.Description)
	if req.Description == constants.Empty {
		zerolog.Ctx(ctx).Warn().Msg("description field is empty, setting to default description")
		return utils.ErrInvalidInput
	}

	return nil
}
