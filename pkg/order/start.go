package order

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	ordercli "github.com/NpoolPlatform/order-middleware/pkg/client/order"

	ordermgrpb "github.com/NpoolPlatform/message/npool/order/mgr/v1/order"
	paymentmgrpb "github.com/NpoolPlatform/message/npool/order/mgr/v1/payment"
	orderpb "github.com/NpoolPlatform/message/npool/order/mw/v1/order"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
)

func orderStart(order *orderpb.Order) (bool, error) {
	switch order.PaymentState {
	case paymentmgrpb.PaymentState_Wait:
		fallthrough // nolint
	case paymentmgrpb.PaymentState_Canceled:
		fallthrough // nolint
	case paymentmgrpb.PaymentState_TimeOut:
		return false, nil
	case paymentmgrpb.PaymentState_Done:
	default:
		return false, fmt.Errorf("invalid payment state")
	}

	if uint32(time.Now().Unix()) < order.Start {
		return false, nil
	}

	return true, nil
}

func processOrderStart(ctx context.Context, order *orderpb.Order) error {
	start, err := orderStart(order)
	if err != nil {
		return err
	}
	if !start {
		return nil
	}

	ostate := ordermgrpb.OrderState_InService
	logger.Sugar().Infow("processOrderStart",
		"OrderID", order.ID,
		"Start", order.Start,
		"State", order.OrderState,
		"NewState", ostate,
		"Units", order.Units,
	)
	_, err = ordercli.UpdateOrder(ctx, &orderpb.OrderReq{
		ID:        &order.ID,
		PaymentID: &order.PaymentID,
		State:     &ostate,
	})

	if err != nil {
		return err
	}

	units, err := decimal.NewFromString(order.Units)
	if err != nil {
		return err
	}
	err = updateStock(ctx, order.GoodID, decimal.NewFromInt(0), units, units.Neg())
	if err != nil {
		return err
	}

	return nil
}

func checkOrderStart(ctx context.Context) {
	offset := int32(0)
	limit := int32(1000)

	for {
		orders, _, err := ordercli.GetOrders(ctx, &orderpb.Conds{
			State: &commonpb.Uint32Val{
				Op:    cruder.EQ,
				Value: uint32(ordermgrpb.OrderState_Paid),
			},
		}, offset, limit)
		if err != nil {
			logger.Sugar().Errorw("processOrders", "offset", offset, "limit", limit, "error", err)
			return
		}
		if len(orders) == 0 {
			return
		}

		for index, order := range orders {
			if err := processOrderStart(ctx, order); err != nil {
				logger.Sugar().Errorw("processOrders", "offset", offset, "index", index, "error", err)
				return
			}
		}

		offset += limit
	}
}
