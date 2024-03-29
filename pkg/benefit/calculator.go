package benefit

import (
	"context"
	"fmt"
	"time"

	timedef "github.com/NpoolPlatform/go-service-framework/pkg/const/time"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"

	gbmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/goodbenefit"
	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/goodbenefit"

	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"
	paymentmgrpb "github.com/NpoolPlatform/message/npool/order/mgr/v1/payment"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"

	ordermgrpb "github.com/NpoolPlatform/message/npool/order/mgr/v1/order"
	ordermwpb "github.com/NpoolPlatform/message/npool/order/mw/v1/order"
	ordermwcli "github.com/NpoolPlatform/order-middleware/pkg/client/order"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	appgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/appgood"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/shopspring/decimal"
)

func benefitTimestamp(timestamp uint32, interval time.Duration) uint32 {
	intervalFloat := interval.Seconds()
	intervalUint := uint32(intervalFloat)
	return timestamp / intervalUint * intervalUint
}

func (st *State) coin(ctx context.Context, coinTypeID string) (*coinmwpb.Coin, error) {
	coin, ok := st.Coins[coinTypeID]
	if ok {
		return coin, nil
	}

	coin, err := coinmwcli.GetCoin(ctx, coinTypeID)
	if err != nil {
		return nil, err
	}

	st.Coins[coinTypeID] = coin

	return coin, nil
}

func (st *State) goodBenefit(ctx context.Context, good *Good) (*gbmwpb.Account, error) {
	acc, ok := st.GoodBenefits[good.ID]
	if ok {
		good.Retry = true
		return acc, nil
	}

	acc, err := gbmwcli.GetAccountOnly(ctx, &gbmwpb.Conds{
		GoodID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: good.ID,
		},
		Backup: &basetypes.BoolVal{
			Op:    cruder.EQ,
			Value: false,
		},
		Active: &basetypes.BoolVal{
			Op:    cruder.EQ,
			Value: true,
		},
		Locked: &basetypes.BoolVal{
			Op:    cruder.EQ,
			Value: false,
		},
		Blocked: &basetypes.BoolVal{
			Op:    cruder.EQ,
			Value: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fail get goodbenefit %v: %v", good.ID, err)
	}
	if acc == nil {
		return nil, fmt.Errorf("invalid goodbenefit %v", good.ID)
	}

	good.Retry = true

	st.GoodBenefits[good.ID] = acc

	return acc, nil
}

func (st *State) balance(ctx context.Context, good *Good) (decimal.Decimal, error) {
	benefit, err := st.goodBenefit(ctx, good)
	if err != nil {
		return decimal.NewFromInt(0), err
	}

	coin, err := st.coin(ctx, good.CoinTypeID)
	if err != nil {
		return decimal.NewFromInt(0), err
	}

	balance, err := sphinxproxycli.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    coin.Name,
		Address: benefit.Address,
	})
	if err != nil {
		return decimal.NewFromInt(0), err
	}

	return decimal.NewFromString(balance.BalanceStr)
}

//nolint:gocognit
func (st *State) CalculateReward(ctx context.Context, good *Good) error {
	total, err := decimal.NewFromString(good.GetGoodTotal())
	if err != nil {
		return err
	}
	if total.Cmp(decimal.NewFromInt(0)) == 0 {
		return fmt.Errorf("invalid stock")
	}

	bal, err := st.balance(ctx, good)
	if err != nil {
		return err
	}

	if bal.Cmp(decimal.NewFromInt(0)) <= 0 {
		return nil
	}

	good.BenefitAccountAmount = bal

	coin, err := st.coin(ctx, good.CoinTypeID)
	if err != nil {
		return err
	}

	reservedAmount, err := decimal.NewFromString(coin.ReservedAmount)
	if err != nil {
		return err
	}

	offset := int32(0)
	limit := int32(100)
	totalInService := decimal.NewFromInt(0)

	for {
		orders, _, err := ordermwcli.GetOrders(ctx, &ordermwpb.Conds{
			GoodID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: good.ID,
			},
			State: &commonpb.Uint32Val{
				Op:    cruder.EQ,
				Value: uint32(ordermgrpb.OrderState_InService),
			},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			break
		}

		for _, ord := range orders {
			units, err := decimal.NewFromString(ord.Units)
			if err != nil {
				return err
			}
			totalInService = totalInService.Add(units)
			if benefitable(good.Good, ord, uint32(time.Now().Unix())) {
				good.BenefitOrderUnits = good.BenefitOrderUnits.Add(units)
			}
		}

		offset += limit
	}

	inService, err := decimal.NewFromString(good.GoodInService)
	if err != nil {
		return err
	}
	if inService.Cmp(totalInService) != 0 {
		logger.Sugar().Errorw("CalculateReward",
			"GoodID", good.ID,
			"GooInService", good.GoodInService,
			"TotalInService", totalInService,
		)
		return fmt.Errorf("inconsistent in service")
	}

	startAmount, _ := decimal.NewFromString(good.NextBenefitStartAmount) //nolint
	good.TodayRewardAmount = bal.
		Sub(reservedAmount).
		Sub(startAmount)
	if good.TodayRewardAmount.Cmp(decimal.NewFromInt(0)) < 0 {
		logger.Sugar().Errorw("CalculateReward",
			"GoodID", good.ID,
			"TodayReward", good.TodayRewardAmount,
			"Balance", bal,
			"StartAmount", startAmount,
			"ReservedAmount", reservedAmount,
		)
		return fmt.Errorf("invalid reward amount")
	}

	good.UserRewardAmount = good.TodayRewardAmount.
		Mul(good.BenefitOrderUnits).
		Div(total)
	good.PlatformRewardAmount = good.TodayRewardAmount.
		Sub(good.UserRewardAmount)

	logger.Sugar().Infow("CalculateReward",
		"GoodID", good.ID,
		"TodayReward", good.TodayRewardAmount,
		"UserReward", good.UserRewardAmount,
		"PlatformReward", good.PlatformRewardAmount,
		"TotalInService", totalInService,
		"BenefitOrderUnits", good.BenefitOrderUnits,
		"StartAmount", startAmount,
		"ReservedAmount", reservedAmount,
		"Balance", bal,
	)

	return nil
}

func benefitable(good *goodmwpb.Good, order *ordermwpb.Order, dateTime uint32) bool {
	if order.PaymentState != paymentmgrpb.PaymentState_Done {
		return false
	}

	orderEnd := order.Start + uint32(good.DurationDays*timedef.SecondsPerDay)
	if orderEnd < dateTime {
		return false
	}
	if order.Start > dateTime {
		return false
	}
	if dateTime < order.Start+uint32(benefitInterval.Seconds()) {
		return false
	}

	return true
}

//nolint:gocognit
func (st *State) CalculateTechniqueServiceFee(ctx context.Context, good *Good) error {
	appUnits := map[string]decimal.Decimal{}
	offset := int32(0)
	limit := int32(100)

	for {
		orders, _, err := ordermwcli.GetOrders(ctx, &ordermwpb.Conds{
			GoodID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: good.ID,
			},
			State: &commonpb.Uint32Val{
				Op:    cruder.EQ,
				Value: uint32(ordermgrpb.OrderState_InService),
			},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			break
		}

		for _, ord := range orders {
			if !benefitable(good.Good, ord, uint32(time.Now().Unix())) {
				continue
			}
			units, err := decimal.NewFromString(ord.Units)
			if err != nil {
				return err
			}
			appUnits[ord.AppID] = appUnits[ord.AppID].Add(units)
			good.BenefitOrderIDs[ord.ID] = ord.PaymentID
		}

		offset += limit
	}

	appIDs := []string{}
	for appID := range appUnits {
		appIDs = append(appIDs, appID)
	}

	goodInService, err := decimal.NewFromString(good.GoodInService)
	if err != nil {
		return err
	}
	if good.BenefitOrderUnits.Cmp(goodInService) > 0 {
		return fmt.Errorf("inconsistent in service")
	}

	goods, _, err := appgoodmwcli.GetGoods(ctx, &appgoodmgrpb.Conds{
		AppIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: appIDs,
		},
		GoodID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: good.ID,
		},
	}, int32(0), int32(len(appIDs)))
	if err != nil {
		return err
	}

	goodMap := map[string]*appgoodmwpb.Good{}
	for _, g := range goods {
		goodMap[g.AppID] = g
	}

	techniqueServiceFee := decimal.NewFromInt(0)

	if good.BenefitOrderUnits.Cmp(decimal.NewFromInt(0)) > 0 {
		for appID, units := range appUnits {
			ag, ok := goodMap[appID]
			if !ok {
				return fmt.Errorf("unauthorized appgood")
			}

			_fee := good.UserRewardAmount.
				Mul(units).
				Div(good.BenefitOrderUnits).
				Mul(decimal.NewFromInt(int64(ag.TechnicalFeeRatio))).
				Div(decimal.NewFromInt(100))

			logger.Sugar().Infow("CalculateTechniqueServiceFee",
				"GoodID", good.ID,
				"GoodName", good.Title,
				"TotalInService", good.GoodInService,
				"BenefitOrderUnits", good.BenefitOrderUnits,
				"AppID", appID,
				"Units", units,
				"Orders", len(good.BenefitOrderIDs),
				"TechnicalFeeRatio", ag.TechnicalFeeRatio,
				"FeeAmount", _fee,
			)

			techniqueServiceFee = techniqueServiceFee.Add(_fee)
		}
	}

	good.TechniqueServiceFeeAmount = techniqueServiceFee
	good.UserRewardAmount = good.UserRewardAmount.Sub(techniqueServiceFee)

	return nil
}
