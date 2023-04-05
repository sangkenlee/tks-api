package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/openinfradev/tks-api/internal/repository"
	"github.com/openinfradev/tks-api/pkg/domain"
)

type IStackTemplateUsecase interface {
	Get(stackTemplate uuid.UUID) (domain.StackTemplate, error)
	Fetch(organizationId string) ([]domain.StackTemplate, error)
	Create(ctx context.Context, dto domain.StackTemplate) (stackTemplate uuid.UUID, err error)
	Update(ctx context.Context, dto domain.StackTemplate) error
	Delete(ctx context.Context, dto domain.StackTemplate) error
}

type StackTemplateUsecase struct {
	repo repository.IStackTemplateRepository
}

func NewStackTemplateUsecase(r repository.IStackTemplateRepository) IStackTemplateUsecase {
	return &StackTemplateUsecase{
		repo: r,
	}
}

func (u *StackTemplateUsecase) Create(ctx context.Context, dto domain.StackTemplate) (stackTemplate uuid.UUID, err error) {
	return uuid.Nil, nil
}

func (u *StackTemplateUsecase) Update(ctx context.Context, dto domain.StackTemplate) error {
	return nil
}

func (u *StackTemplateUsecase) Get(stackTemplate uuid.UUID) (res domain.StackTemplate, err error) {
	res, err = u.repo.Get(stackTemplate)
	if err != nil {
		return domain.StackTemplate{}, err
	}
	return
}

func (u *StackTemplateUsecase) Fetch(organizationId string) (res []domain.StackTemplate, err error) {
	res, err = u.repo.Fetch(organizationId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *StackTemplateUsecase) Delete(ctx context.Context, dto domain.StackTemplate) (err error) {
	return nil
}
