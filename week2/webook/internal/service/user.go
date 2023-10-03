// Copyright@daidai53 2023
package service

import (
	"context"
	"errors"
	"github.com/daidai53/go-homework/week2/webook/internal/domain"
	"github.com/daidai53/go-homework/week2/webook/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户邮箱或者密码不存在")
	ErrNoUserProfile         = repository.ErrUserNotFound
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (u *UserService) SignUp(ctx context.Context, user domain.User) error {
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(encryptedPwd)
	return u.repo.Create(ctx, user)
}

func (u *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	usr, err := u.repo.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return usr, nil
}

func (u *UserService) Profile(c *gin.Context, userId int64) error {
	usr, err := u.repo.FindById(c, userId)
	if errors.Is(err, repository.ErrUserNotFound) {
		c.String(http.StatusOK, "系统错误")
		return err
	}
	c.JSON(http.StatusOK, toJson(usr))
	return nil
}

func (u *UserService) Edit(c *gin.Context, idInt64 int64, nickname, birthday, aboutMe string) error {
	return u.repo.Update(c, idInt64, domain.User{
		Nickname: nickname,
		Birthday: birthday,
		AboutMe:  aboutMe,
	})
}

func toJson(usr domain.User) gin.H {
	return gin.H{
		"Email":    usr.Email,
		"Nickname": usr.Nickname,
		"Phone":    usr.Phone,
		"Birthday": usr.Birthday,
		"AboutMe":  usr.AboutMe,
	}
}
