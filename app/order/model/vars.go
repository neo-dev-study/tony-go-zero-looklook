package model

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrNotFound     = sqlx.ErrNotFound
	ErrNoRowsUpdate = errors.New("update db no rows change")
)

// HomestayOrder 交易状态 :  -1: 已取消 0:待支付 1:未使用 2:已使用  3:已过期

var (
	HomestayOrderTradeStateCancel  int64 = -1
	HomestayOrderTradeStateWaitPay int64 = 0
	HomestayOrderTradeStateWaitUse int64 = 1
	HomestayOrderTradeStateUsed    int64 = 2
	HomestayOrderTradeStateRefund  int64 = 3
	HomestayOrderTradeStateExpire  int64 = 4
)

// 是否需要餐食

var (
	HomestayOrderNeedFoodNo  int64 = 0
	HomestayOrderNeedFoodYes int64 = 1
)
