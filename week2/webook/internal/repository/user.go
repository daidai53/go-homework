// Copyright@daidai53 2023
package repository

import (
	"context"
	"github.com/daidai53/go-homework/week2/webook/internal/domain"
	"github.com/daidai53/go-homework/week2/webook/internal/repository/dao"
	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{dao: dao}
}

func (u *UserRepository) Create(ctx context.Context, user domain.User) error {
	return u.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	usr, err := u.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return u.toDomainUser(usr), nil
}

func (u *UserRepository) FindById(ctx *gin.Context, id int64) (domain.User, error) {
	usr, err := u.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return u.toDomainUser(usr), nil
}

func (u *UserRepository) toDomainUser(usr dao.User) domain.User {
	return domain.User{
		Id:       usr.Id,
		Email:    usr.Email,
		Password: usr.Password,
		Nickname: usr.Nickname,
		Phone:    usr.Phone,
		Birthday: usr.Birthday,
		AboutMe:  usr.AboutMe,
	}
}

func (u *UserRepository) Update(ctx context.Context, idInt64 int64, user domain.User) error {
	return u.dao.Update(ctx, idInt64, dao.User{
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		AboutMe:  user.AboutMe,
	})
}
