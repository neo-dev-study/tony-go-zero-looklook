package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"looklook/app/travel/cmd/rpc/internal/config"
	"looklook/app/travel/model"
)

type ServiceContext struct {
	Config config.Config

	HomestayModel model.HomestayModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config: c,

		HomestayModel: model.NewHomestayModel(sqlConn, c.Cache),
	}
}
