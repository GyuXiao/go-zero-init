package userTools

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-init/app/admin/models"
	"go-zero-init/common/constant"
	"go-zero-init/common/xerr"
	"strconv"
)

// 校验用户角色是否为管理员

func CheckUserIsAdminRole(client *redis.Redis, token string) error {
	tokenLogic := models.NewDefaultTokenModel(client)
	result, err := tokenLogic.CheckTokenExist(token)
	if err != nil {
		return err
	}

	userRole, _ := strconv.Atoi(result[1])
	if userRole != constant.AdminRole {
		return xerr.NewErrCode(xerr.PermissionDenied)
	}

	return nil
}
