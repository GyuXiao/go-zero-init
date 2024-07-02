package user

import (
	"go-zero-init/common/result"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-init/app/user/cmd/api/internal/logic/user"
	"go-zero-init/app/user/cmd/api/internal/svc"
	"go-zero-init/app/user/cmd/api/internal/types"
)

func CurrentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CurrentUserReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewCurrentLogic(r.Context(), svcCtx)
		resp, err := l.Current(&req)
		result.HttpResult(r, w, resp, err)
	}
}
