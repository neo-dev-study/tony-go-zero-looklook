package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type OrderConfig struct {
	zrpc.RpcServerConf

	DB struct {
		DataSource string
	}
	Cache cache.CacheConf

	TravelRpcConf zrpc.RpcClientConf
}
