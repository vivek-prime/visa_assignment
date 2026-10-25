package service

import (
	"context"
	"errors"
	"testing"
	"time"
	dto "visa_assignment/models/dto"
	entity "visa_assignment/models/entity"
	"visa_assignment/utils"
)

func TestCreateMember_Success(t *testing.T) {
	memberRepo := &memberRepoMock{}
	rewardRepo := &rewardRepoMock{}

	svc := NewMemberService(memberRepo, rewardRepo)

	member, err := svc.CreateMember(context.Background(), dto.MemberRequest{
		Name:  "  Alex  ",
		Email: " ALEX@example.com ",
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if member.Name != "Alex" {
		t.Errorf("expected name Alex, got %q", member.Name)
	}

	if member.Email != "alex@example.com" {
		t.Errorf("expected lowercase email, got %q", member.Email)
	}

	if member.MemberID != 1 {
		t.Errorf("expected member ID 1, got %d", member.MemberID)
	}

	if member.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestCreateMember_EmptyName(t *testing.T) {
	svc := NewMemberService(&memberRepoMock{}, &rewardRepoMock{})

	member, err := svc.CreateMember(context.Background(), dto.MemberRequest{
		Name:  " ",
		Email: "alex@example.com",
	})

	if member != nil {
		t.Error("expected nil member")
	}

	if !errors.Is(err, utils.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreateMember_EmptyEmail(t *testing.T) {
	svc := NewMemberService(&memberRepoMock{}, &rewardRepoMock{})

	member, err := svc.CreateMember(context.Background(), dto.MemberRequest{
		Name:  "Alex",
		Email: " ",
	})

	if member != nil {
		t.Error("expected nil member")
	}

	if !errors.Is(err, utils.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCreateMember_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	memberRepo := &memberRepoMock{err: expectedErr}
	svc := NewMemberService(memberRepo, &rewardRepoMock{})

	member, err := svc.CreateMember(context.Background(), dto.MemberRequest{
		Name:  "Alex",
		Email: "alex@example.com",
	})

	if member != nil {
		t.Error("expected nil member")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected database error, got %v", err)
	}
}

func TestGetMemberDetails_Success(t *testing.T) {
	createdAt := time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC)

	memberRepo := &memberRepoMock{
		member: &entity.MemberEntity{
			MemberID:  1,
			Name:      "Alex",
			Email:     "alex@example.com",
			CreatedAt: createdAt,
		},
	}

	rewardRepo := &rewardRepoMock{balance: 500}
	svc := NewMemberService(memberRepo, rewardRepo)

	result, err := svc.GetMemberDetails(context.Background(), 1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if result.MemberID != 1 {
		t.Errorf("expected member ID 1, got %d", result.MemberID)
	}

	if result.Name != "Alex" {
		t.Errorf("expected name Alex, got %q", result.Name)
	}

	if result.Email != "alex@example.com" {
		t.Errorf("expected email alex@example.com, got %q", result.Email)
	}

	if result.PointsBalance != 500 {
		t.Errorf("expected balance 500, got %v", result.PointsBalance)
	}

	if !result.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt %v, got %v", createdAt, result.CreatedAt)
	}
}

func TestGetMemberDetails_MemberNotFound(t *testing.T) {
	svc := NewMemberService(&memberRepoMock{}, &rewardRepoMock{})

	result, err := svc.GetMemberDetails(context.Background(), 99)

	if result != nil {
		t.Error("expected nil result")
	}

	if !errors.Is(err, utils.ErrMemberNotFound) {
		t.Errorf("expected ErrMemberNotFound, got %v", err)
	}
}

func TestGetMemberDetails_RewardRepositoryError(t *testing.T) {
	expectedErr := errors.New("reward database error")

	memberRepo := &memberRepoMock{
		member: &entity.MemberEntity{
			MemberID: 1,
			Name:     "Alex",
			Email:    "alex@example.com",
		},
	}

	rewardRepo := &rewardRepoMock{err: expectedErr}
	svc := NewMemberService(memberRepo, rewardRepo)

	result, err := svc.GetMemberDetails(context.Background(), 1)

	if result != nil {
		t.Error("expected nil result")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected reward database error, got %v", err)
	}
}
