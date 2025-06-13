package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type PaymentConfig struct {
	zrpc.RpcServerConf

	DB struct {
		DataSource string
	}
	Cache                        cache.CacheConf
	KqPaymentUpdatePayStatusConf KqConfig
}
