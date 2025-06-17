package svc

import (
	"github.com/zeromicro/go-zero/zrpc"

	"looklook/app/order/cmd/api/internal/config"
	"looklook/app/order/cmd/rpc/order"
	"looklook/app/payment/cmd/rpc/payment"
	"looklook/app/travel/cmd/rpc/travel"
)

type ServiceContext struct {
	Config config.OrderConfig

	OrderRpc   order.Order
	PaymentRpc payment.Payment
	TravelRpc  travel.Travel
}

func NewServiceContext(c config.OrderConfig) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:   order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		PaymentRpc: payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		TravelRpc:  travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
