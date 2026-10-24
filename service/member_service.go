package service

import (
	"context"
	"strings"
	"time"

	dto "visa_assignment/models/dto"
	entity "visa_assignment/models/entity"
	"visa_assignment/repository"
	"visa_assignment/utils"
)

type MemberService interface {
	CreateMember(ctx context.Context, req dto.MemberRequest) (*entity.MemberEntity, error)
	GetMemberWithBalance(ctx context.Context, memberID int64) (*dto.MemberBalanceResponse, error)
}

type memberService struct {
	memberRepo repository.MemberRepository
	rewardRepo repository.RewardRepository
}

func NewMemberService(memberRepo repository.MemberRepository, rewardRepo repository.RewardRepository) MemberService {
	return &memberService{
		memberRepo: memberRepo,
		rewardRepo: rewardRepo,
	}
}

func (s *memberService) CreateMember(ctx context.Context, req dto.MemberRequest) (*entity.MemberEntity, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if name == "" || email == "" {
		return nil, utils.ErrInvalidInput
	}

	// Rule: A member cannot be created with a duplicate email
	exists, err := s.memberRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.ErrDuplicateEmail
	}

	member := &entity.MemberEntity{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *memberService) GetMemberWithBalance(ctx context.Context, memberID int64) (*dto.MemberBalanceResponse, error) {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return nil, err
	}

	// Rule: A member's available balance is the sum of all their reward entries
	balance, err := s.rewardRepo.GetMemberBalance(ctx, memberID)
	if err != nil {
		return nil, err
	}

	return &dto.MemberBalanceResponse{
		MemberID:      member.MemberID,
		Name:          member.Name,
		Email:         member.Email,
		PointsBalance: balance,
		CreatedAt:     member.CreatedAt,
	}, nil
}
