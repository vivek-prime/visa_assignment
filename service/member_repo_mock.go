package service

import (
	"context"
	entity "visa_assignment/models/entity"
	"visa_assignment/utils"
)

type memberRepoMock struct {
	member *entity.MemberEntity
	err    error
}

func (m *memberRepoMock) Create(ctx context.Context, member *entity.MemberEntity) error {
	if m.err != nil {
		return m.err
	}

	member.MemberID = 1
	m.member = member
	return nil
}

func (m *memberRepoMock) GetByID(ctx context.Context, memberID int64) (*entity.MemberEntity, error) {
	if m.err != nil {
		return nil, m.err
	}

	if m.member == nil {
		return nil, utils.ErrMemberNotFound
	}

	return m.member, nil
}
