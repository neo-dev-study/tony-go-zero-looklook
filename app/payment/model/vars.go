package model

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrNotFound     = sqlx.ErrNotFound
	ErrNoRowsUpdate = errors.New("update db no rows change")
)

// 支付业务类型

var ThirdPaymentServiceTypeHomestayOrder string = "homestayOrder" // 民宿支付

// 支付方式

var ThirdPaymentPayModelWechatPay = "WECHAT_PAY" // 微信支付

// 支付状态

var (
	ThirdPaymentPayTradeStateFAIL    int64 = -1 // 支付失败
	ThirdPaymentPayTradeStateWait    int64 = 0  // 待支付
	ThirdPaymentPayTradeStateSuccess int64 = 1  // 支付成功
	ThirdPaymentPayTradeStateRefund  int64 = 2  // 已退款
)
