package logic

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"looklook/app/mqueue/cmd/job/jobtype"
	"looklook/app/order/cmd/rpc/internal/svc"
	"looklook/app/order/cmd/rpc/pb"
	"looklook/app/order/model"
	"looklook/app/travel/cmd/rpc/travel"
	"looklook/common/tool"
	"looklook/common/uniqueid"
	"looklook/common/xerr"
)

const CloseOrderTimeMinutes = 30 // defer close order time

type CreateHomestayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateHomestayOrder.
func (l *CreateHomestayOrderLogic) CreateHomestayOrder(in *pb.CreateHomestayOrderReq) (*pb.CreateHomestayOrderResp, error) {
	// 1、Create Order
	if in.LiveEndTime <= in.LiveStartTime {
		return nil, errors.Wrapf(xerr.NewErrMsg("Stay at least one night"), "Place an order at a B&B. The end time of your stay must be greater than the start time. in : %+v", in)
	}

	homestayDetail, err := l.getHomestayDetail(in.HomestayId)
	if err != nil {
		return nil, err
	}

	order := l.generateOrder(in, homestayDetail)
	_, err = l.svcCtx.HomestayOrderModel.Insert(l.ctx, nil, order)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Order Database Exception order : %+v , err: %v", order, err)
	}

	// 2、Delayed closing of order tasks.
	l.delayClosingOrder(order.Sn)

	return &pb.CreateHomestayOrderResp{
		Sn: order.Sn,
	}, nil
}

func (l *CreateHomestayOrderLogic) getHomestayDetail(homestayId int64) (*travel.HomestayDetailResp, error) {
	resp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &travel.HomestayDetailReq{
		Id: homestayId,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Failed to query the record"), "Failed to query the record  rpc HomestayDetail fail , homestayId : %d , err : %v", homestayId, err)
	}
	if resp.Homestay == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("This record does not exist"), "This record does not exist , homestayId : %d ", homestayId)
	}
	return resp, err
}

func (l *CreateHomestayOrderLogic) generateOrder(in *pb.CreateHomestayOrderReq, homestayDetail *travel.HomestayDetailResp) *model.HomestayOrder {
	var cover string // Get the cover...
	if len(homestayDetail.Homestay.Banner) > 0 {
		cover = strings.Split(homestayDetail.Homestay.Banner, ",")[0]
	}

	order := new(model.HomestayOrder)
	order.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_HOMESTAY_ORDER)
	order.UserId = in.UserId
	order.HomestayId = in.HomestayId
	order.Title = homestayDetail.Homestay.Title
	order.SubTitle = homestayDetail.Homestay.SubTitle
	order.Cover = cover
	order.Info = homestayDetail.Homestay.Info
	order.PeopleNum = homestayDetail.Homestay.PeopleNum
	order.RowType = homestayDetail.Homestay.RowType
	order.HomestayPrice = homestayDetail.Homestay.HomestayPrice
	order.MarketHomestayPrice = homestayDetail.Homestay.MarketHomestayPrice
	order.HomestayBusinessId = homestayDetail.Homestay.HomestayBusinessId
	order.HomestayUserId = homestayDetail.Homestay.UserId
	order.LivePeopleNum = in.LivePeopleNum
	order.TradeState = model.HomestayOrderTradeStateWaitPay
	order.TradeCode = tool.Krand(8, tool.KC_RAND_KIND_ALL)
	order.Remark = in.Remark
	order.FoodInfo = homestayDetail.Homestay.FoodInfo
	order.FoodPrice = homestayDetail.Homestay.FoodPrice
	order.LiveStartDate = time.Unix(in.LiveStartTime, 0)
	order.LiveEndDate = time.Unix(in.LiveEndTime, 0)

	liveDays := int64(order.LiveEndDate.Sub(order.LiveStartDate).Seconds() / 86400) // Stayed a few days in total

	order.HomestayTotalPrice = homestayDetail.Homestay.HomestayPrice * liveDays // Calculate the total price of the B&B
	if in.IsFood {
		order.NeedFood = model.HomestayOrderNeedFoodYes
		// Calculate the total price of the meal.
		order.FoodTotalPrice = homestayDetail.Homestay.FoodPrice * in.LivePeopleNum * liveDays
	}

	order.OrderTotalPrice = order.HomestayTotalPrice + order.FoodTotalPrice // Calculate total order price.
	return order
}

func (l *CreateHomestayOrderLogic) delayClosingOrder(orderSn string) {
	payload, err := json.Marshal(jobtype.DeferCloseHomestayOrderPayload{Sn: orderSn})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("create defer close order task json Marshal fail err :%+v , sn : %s", err, orderSn)
	} else {
		_, err = l.svcCtx.AsynqClient.Enqueue(asynq.NewTask(jobtype.DeferCloseHomestayOrder, payload), asynq.ProcessIn(CloseOrderTimeMinutes*time.Minute))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("create defer close order task insert queue fail err :%+v , sn : %s", err, orderSn)
		}
	}
}
