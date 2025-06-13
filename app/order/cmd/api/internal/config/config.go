package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type OrderConfig struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
	}

	OrderRpcConf   zrpc.RpcClientConf
	PaymentRpcConf zrpc.RpcClientConf
	TravelRpcConf  zrpc.RpcClientConf
}
