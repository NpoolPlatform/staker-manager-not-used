package benefit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	pltfaccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/platform"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	pltfaccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"

	goodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/good"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"

	ordermwpb "github.com/NpoolPlatform/message/npool/order/mw/v1/order"
	ordermwcli "github.com/NpoolPlatform/order-middleware/pkg/client/order"

	txmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/tx"
	txmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/tx"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	notifbenefitpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/goodbenefit"
	notifbenefitcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/goodbenefit"
	"github.com/shopspring/decimal"
)

func (st *State) platformAccount(
	ctx context.Context,
	coinTypeID string,
	usedFor basetypes.AccountUsedFor,
) (
	*pltfaccmwpb.Account,
	error,
) {
	accs, ok := st.PlatformAccounts[coinTypeID]
	if ok {
		acc, ok := accs[usedFor]
		if ok {
			return acc, nil
		}
	}

	acc, err := pltfaccmwcli.GetAccountOnly(ctx, &pltfaccmwpb.Conds{
		CoinTypeID: &basetypes.StringVal{
			Op:    cruder.EQ,
			Value: coinTypeID,
		},
		UsedFor: &basetypes.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(usedFor),
		},
		Backup: &basetypes.BoolVal{
			Op:    cruder.EQ,
			Value: false,
		},
		Active: &basetypes.BoolVal{
			Op:    cruder.EQ,
			Value: true,
		},
		Blocked: &basetypes.BoolVal{
			Op:    cruder.EQ,
			Value: false,
		},
		Locked: &basetypes.BoolVal{
			Op:    cruder.EQ,
			Value: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fail get platform %v usedfor %v: %v", coinTypeID, usedFor, err)
	}
	if acc == nil {
		return nil, fmt.Errorf("invalid account of cointypeid %v usedfor %v", coinTypeID, usedFor)
	}

	_, ok = st.PlatformAccounts[coinTypeID]
	if !ok {
		st.PlatformAccounts[coinTypeID] = map[basetypes.AccountUsedFor]*pltfaccmwpb.Account{}
	}
	st.PlatformAccounts[coinTypeID][usedFor] = acc

	return acc, nil
}

func (st *State) TransferReward(ctx context.Context, good *Good) error { //nolint
	if good.TodayRewardAmount.Cmp(decimal.NewFromInt(0)) <= 0 {
		return nil
	}

	userHotAcc, err := st.platformAccount(
		ctx,
		good.CoinTypeID,
		basetypes.AccountUsedFor_UserBenefitHot)
	if err != nil {
		return err
	}

	goodBenefitAcc, err := st.goodBenefit(ctx, good)
	if err != nil {
		return err
	}

	coin, err := coinmwcli.GetCoin(ctx, good.CoinTypeID)
	if err != nil {
		return err
	}
	if coin == nil {
		return fmt.Errorf("invalid coin")
	}

	leastTransferAmount, err := decimal.NewFromString(coin.LeastTransferAmount)
	if err != nil {
		return err
	}
	if leastTransferAmount.Cmp(decimal.NewFromInt(0)) <= 0 {
		return fmt.Errorf("invalid least transfer amount")
	}

	txs := []*txmwpb.TxReq{}
	toUser := decimal.NewFromInt(0)

	if good.TodayRewardAmount.Cmp(leastTransferAmount) > 0 {
		toUser = good.TodayRewardAmount
	}

	if toUser.Cmp(decimal.NewFromInt(0)) > 0 {
		amount := toUser.String()
		feeAmount := decimal.NewFromInt(0).String()
		txExtra := fmt.Sprintf(
			`{"GoodID":"%v","Reward":"%v","UserReward":"%v","PlatformReward":"%v","TechniqueServiceFee":"%v"}`,
			good.ID,
			good.TodayRewardAmount,
			good.UserRewardAmount,
			good.PlatformRewardAmount,
			good.TechniqueServiceFeeAmount,
		)
		txType := basetypes.TxType_TxUserBenefit
		txs = append(txs, &txmwpb.TxReq{
			CoinTypeID:    &good.CoinTypeID,
			FromAccountID: &goodBenefitAcc.AccountID,
			ToAccountID:   &userHotAcc.AccountID,
			Amount:        &amount,
			FeeAmount:     &feeAmount,
			Extra:         &txExtra,
			Type:          &txType,
		})
	}

	reservedAmount, err := decimal.NewFromString(coin.ReservedAmount)
	if err != nil {
		return err
	}

	nextStartAmount := good.BenefitAccountAmount.Sub(reservedAmount)
	state := goodmgrpb.BenefitState_BenefitTransferring
	nextStartAmountS := nextStartAmount.String()
	lastBenefitAmountS := good.TodayRewardAmount.String()

	req := &goodmwpb.GoodReq{
		ID:                     &good.ID,
		BenefitState:           &state,
		NextBenefitStartAmount: &nextStartAmountS,
		LastBenefitAmount:      &lastBenefitAmountS,
	}
	g, err := goodmwcli.UpdateGood(ctx, req)
	if err != nil {
		return err
	}

	if good.UserRewardAmount.Cmp(decimal.NewFromInt(0)) > 0 {
		ords := []*ordermwpb.OrderReq{}
		for _oid, _pid := range good.BenefitOrderIDs {
			oid := _oid
			pid := _pid
			ords = append(ords, &ordermwpb.OrderReq{
				ID:            &oid,
				PaymentID:     &pid,
				LastBenefitAt: &g.LastBenefitAt,
			})
		}
		if len(ords) > 0 {
			logger.Sugar().Infow("TransferReward",
				"GoodID", good.ID,
				"UserRewardAmount", good.UserRewardAmount,
				"Units", good.BenefitOrderUnits,
				"Orders", len(good.BenefitOrderIDs),
				"UpdateOrders", len(ords),
				"LastBenefitAt", g.LastBenefitAt,
			)
			_, err := ordermwcli.UpdateOrders(ctx, ords)
			if err != nil {
				return err
			}
		}
	}

	if len(txs) == 0 {
		return nil
	}

	infos, err := txmwcli.CreateTxs(ctx, txs)
	if err != nil {
		return err
	}

	req.BenefitState = nil
	for _, tx := range infos {
		req.BenefitTIDs = append(req.BenefitTIDs, tx.ID)
	}
	_, err = goodmwcli.UpdateGood(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

//nolint
func (st *State) CheckTransfer(ctx context.Context, good *Good) error {
	transferred := decimal.NewFromInt(0)
	toPlatform := decimal.NewFromInt(0)
	doneTIDs := []string{}
	txFail := false
	txExtra := ""

	logger.Sugar().Infow(
		"CheckTransfer",
		"GoodID", good.ID,
		"BenefitTIDs", good.BenefitTIDs,
	)

	if len(good.BenefitTIDs) > 0 {
		txs, _, err := txmwcli.GetTxs(ctx, &txmwpb.Conds{
			IDs: &basetypes.StringSliceVal{
				Op:    cruder.IN,
				Value: good.BenefitTIDs,
			},
		}, int32(0), int32(len(good.BenefitTIDs)))
		if err != nil {
			return err
		}

		for _, tx := range txs {
			switch tx.Type {
			case basetypes.TxType_TxPlatformBenefit:
			case basetypes.TxType_TxUserBenefit:
			default:
				return fmt.Errorf("invalid tx type")
			}

			switch tx.State {
			case basetypes.TxState_TxStateCreated:
				fallthrough //nolint
			case basetypes.TxState_TxStateWait:
				fallthrough //nolint
			case basetypes.TxState_TxStateTransferring:
				return nil
			case basetypes.TxState_TxStateFail:
				txFail = true
				fallthrough //nolint
			case basetypes.TxState_TxStateSuccessful:
				amount, err := decimal.NewFromString(tx.Amount)
				if err != nil {
					return err
				}
				transferred = transferred.Add(amount)
				doneTIDs = append(doneTIDs, tx.ID)

				type p struct {
					PlatformReward      decimal.Decimal
					TechniqueServiceFee decimal.Decimal
					GoodID              string
				}
				_p := p{}
				err = json.Unmarshal([]byte(tx.Extra), &_p)
				if err != nil {
					return err
				}

				toPlatform = _p.PlatformReward.Add(_p.TechniqueServiceFee)
				txExtra = tx.Extra

				now := uint32(time.Now().Unix())
				_result := basetypes.Result(basetypes.Result_value[basetypes.Result_Fail.String()])
				if tx.State == basetypes.TxState_TxStateSuccessful {
					_result = basetypes.Result(basetypes.Result_value[basetypes.Result_Success.String()])
				}

				message := tx.State.String()
				_, err = notifbenefitcli.CreateGoodBenefit(ctx, &notifbenefitpb.GoodBenefitReq{
					GoodID:      &_p.GoodID,
					GoodName:    &good.Title,
					Amount:      &tx.Amount,
					TxID:        &tx.ID,
					Message:     &message,
					State:       &_result,
					BenefitDate: &now,
				})
				if err != nil {
					logger.Sugar().Errorw("CreateGoodBenefit", "Error", err)
				}
			}
		}
	}

	// TODO: we need to improve for some tx fail, some tx success

	nextStart, err := decimal.NewFromString(good.NextBenefitStartAmount)
	if err != nil {
		return err
	}

	nextStart = nextStart.Sub(transferred)
	if nextStart.Cmp(decimal.NewFromInt(0)) < 0 {
		return fmt.Errorf("invalid start amount nextStart %v, transferred %v", nextStart, transferred)
	}

	logger.Sugar().Errorw(
		"CheckTransfer",
		"GoodID", good.ID,
		"Transferred", transferred,
		"NextStart", nextStart,
		"BenefitTIDs", good.BenefitTIDs,
		"DoneTIDs", doneTIDs,
	)

	state := goodmgrpb.BenefitState_BenefitBookKeeping
	nextStartS := nextStart.String()

	remainTIDs := []string{}

	for _, tid := range good.BenefitTIDs {
		found := false
		for _, _tid := range doneTIDs {
			if _tid == tid {
				found = true
				break
			}
		}
		if !found {
			remainTIDs = append(remainTIDs, tid)
		}
	}

	req := &goodmwpb.GoodReq{
		ID:                     &good.ID,
		BenefitState:           &state,
		NextBenefitStartAmount: &nextStartS,
		BenefitTIDs:            remainTIDs,
	}

	_, err = goodmwcli.UpdateGood(ctx, req)
	if err != nil {
		return err
	}

	if !txFail {
		coin, err := coinmwcli.GetCoin(ctx, good.CoinTypeID)
		if err != nil {
			logger.Sugar().Warnw(
				"CheckTransfer",
				"Extra", txExtra,
				"Error", err,
			)
			return nil
		}
		if coin == nil {
			logger.Sugar().Warnw(
				"CheckTransfer",
				"Extra", txExtra,
				"Error", "invalid coin",
			)
			return nil
		}

		leastTransferAmount, err := decimal.NewFromString(coin.LeastTransferAmount)
		if err != nil {
			logger.Sugar().Warnw(
				"CheckTransfer",
				"Extra", txExtra,
				"Error", err,
			)
			return nil
		}
		if leastTransferAmount.Cmp(decimal.NewFromInt(0)) <= 0 {
			logger.Sugar().Warnw(
				"CheckTransfer",
				"Extra", txExtra,
				"Error", "invalid least transfer amount",
			)
			return nil
		}

		if toPlatform.Cmp(leastTransferAmount) > 0 {
			userHotAcc, err := st.platformAccount(
				ctx,
				good.CoinTypeID,
				basetypes.AccountUsedFor_UserBenefitHot)
			if err != nil {
				logger.Sugar().Warnw(
					"CheckTransfer",
					"Extra", txExtra,
					"Error", err,
				)
				return nil
			}

			pltfColdAcc, err := st.platformAccount(
				ctx,
				good.CoinTypeID,
				basetypes.AccountUsedFor_PlatformBenefitCold)
			if err != nil {
				logger.Sugar().Warnw(
					"CheckTransfer",
					"Extra", txExtra,
					"Error", err,
				)
				return nil
			}

			amount := toPlatform.String()
			feeAmount := decimal.NewFromInt(0).String()
			txType := basetypes.TxType_TxPlatformBenefit
			_, err = txmwcli.CreateTx(ctx, &txmwpb.TxReq{
				CoinTypeID:    &good.CoinTypeID,
				FromAccountID: &userHotAcc.AccountID,
				ToAccountID:   &pltfColdAcc.AccountID,
				Amount:        &amount,
				FeeAmount:     &feeAmount,
				Extra:         &txExtra,
				Type:          &txType,
			})
			if err != nil {
				logger.Sugar().Warnw(
					"CheckTransfer",
					"Extra", txExtra,
					"Error", err,
				)
			}
		}
		return nil
	}

	state = goodmgrpb.BenefitState_BenefitWait
	req = &goodmwpb.GoodReq{
		ID:           &good.ID,
		BenefitState: &state,
	}

	_, err = goodmwcli.UpdateGood(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
