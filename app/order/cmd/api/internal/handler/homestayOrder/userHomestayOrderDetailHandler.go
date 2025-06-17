package homestayOrder

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"looklook/app/order/cmd/api/internal/logic/homestayOrder"
	"looklook/app/order/cmd/api/internal/svc"
	"looklook/app/order/cmd/api/internal/types"
	"looklook/common/result"
)

func UserHomestayOrderDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserHomestayOrderDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestayOrder.NewUserHomestayOrderDetailLogic(r.Context(), ctx)
		resp, err := l.UserHomestayOrderDetail(req)
		result.HttpResult(r, w, resp, err)
	}
}
