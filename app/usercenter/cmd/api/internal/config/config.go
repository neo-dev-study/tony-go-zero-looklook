package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type UsercenterConfig struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
	}
	WxMiniConf        WxMiniConf
	UsercenterRpcConf zrpc.RpcClientConf
}
