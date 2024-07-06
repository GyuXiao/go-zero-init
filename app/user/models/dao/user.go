package dao

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"go-zero-init/app/user/models/do"
	"go-zero-init/app/user/models/entity"
	"go-zero-init/common/xerr"
	"gorm.io/gorm"
	"sync"
)

type UserService interface {
	CreateUser(*entity.User) error
	SearchUserByUsername(string) (*entity.User, error)
}

var userService UserService
var userOnce sync.Once

type defaultUserModel struct {
	*do.Query
}

func NewDefaultUserModel(db *do.Query) UserService {
	userOnce.Do(func() {
		userService = &defaultUserModel{db}
	})
	return userService
}

func (m *defaultUserModel) CreateUser(userMap *entity.User) error {
	err := m.User.Create(userMap)
	if err != nil {
		logc.Infof(ctx, "mysql create user err: %v", err)
		return xerr.NewErrCode(xerr.CreateUserError)
	}
	return nil
}

func (m *defaultUserModel) SearchUserByUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	user, err := m.User.Where(m.User.Username.Value(username)).Where(m.User.IsDelete.Value(0)).First()
	switch {
	case err == nil:
		return user, nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		logc.Info(ctx, "mysql search user by username not found")
		return nil, xerr.NewErrCode(xerr.RecordNotFoundError)
	default:
		logc.Infof(ctx, "mysql search user by username err: %v", err)
		return nil, xerr.NewErrCode(xerr.SearchUserError)
	}
}
