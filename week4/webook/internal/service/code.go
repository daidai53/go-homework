// Copyright@daidai53 2023
package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/daidai53/webook/internal/repository"
	"github.com/daidai53/webook/internal/service/sms"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type codeService struct {
	repo repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(c repository.CodeRepository, sms sms.Service) CodeService {
	return &codeService{
		repo: c,
		sms:  sms,
	}
}

func (c *codeService) Send(ctx context.Context, biz, phone string) error {
	code := c.generateCode()
	err := c.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return c.sms.Send(ctx, codeTplId, []string{code}, phone)
}

func (c *codeService) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	verify, err := c.repo.Verify(ctx, biz, phone, code)
	if errors.Is(err, repository.ErrCodeSendTooMany) {
		return false, nil
	}
	return verify, err
}

func (c *codeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
