// Copyright@daidai53 2023
package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

type User struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Email      string `gorm:"unique"`
	Password   string
	CreateTime int64
	UpdateTime int64

	Nickname string
	Phone    string
	Birthday string
	AboutMe  string
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			const duplicateErr uint16 = 1062
			if mysqlErr.Number == duplicateErr {
				return ErrDuplicateEmail
			}
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&u).Error
	return u, err
}

func (dao *UserDAO) Update(ctx context.Context, idInt64 int64, user User) error {
	return dao.db.WithContext(ctx).Where("id=?", idInt64).Updates(&user).Error
}
