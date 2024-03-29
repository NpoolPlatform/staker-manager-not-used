package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	useraccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/user"
	txmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/tx"
	notifmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"
	txnotifmgrpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/tx"
	tmplmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template"

	useraccmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/user"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	txmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/tx"
	notifmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"
	txnotifmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/tx"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func waitSuccess(ctx context.Context) error { //nolint
	offset := int32(0)
	limit := int32(1000)

	for {
		notifs, _, err := txnotifmwcli.GetTxs(ctx, &txnotifmgrpb.Conds{
			NotifState: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(txnotifmgrpb.TxState_WaitSuccess)},
			TxType:     &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(basetypes.TxType_TxWithdraw)},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(notifs) == 0 {
			break
		}

		tids := []string{}
		notifMap := map[string]*txnotifmgrpb.Tx{}

		for _, notif := range notifs {
			tids = append(tids, notif.TxID)
			notifMap[notif.TxID] = notif
		}

		txs, _, err := txmwcli.GetTxs(ctx, &txmwpb.Conds{
			IDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: tids},
		}, 0, int32(len(tids)))
		if err != nil {
			return err
		}

		for _, tx := range txs {
			if tx.State != basetypes.TxState_TxStateSuccessful {
				continue
			}

			acc, err := useraccmwcli.GetAccountOnly(ctx, &useraccmwpb.Conds{
				AccountID: &basetypes.StringVal{Op: cruder.EQ, Value: tx.ToAccountID},
			})
			if err != nil {
				return err
			}
			if acc == nil {
				continue
			}

			user, err := usermwcli.GetUser(ctx, acc.AppID, acc.UserID)
			if err != nil {
				return err
			}
			if user == nil {
				continue
			}

			extra := fmt.Sprintf(`{"TxID":"%v"}`, tx.ID)
			now := uint32(time.Now().Unix())

			notifVars := &tmplmwpb.TemplateVars{
				Username:  &user.Username,
				Amount:    &tx.Amount,
				CoinUnit:  &tx.CoinUnit,
				Address:   &acc.Address,
				Timestamp: &now,
			}

			if _, err := notifmwcli.GenerateNotifs(ctx, &notifmwpb.GenerateNotifsRequest{
				AppID:     acc.AppID,
				UserID:    acc.UserID,
				EventType: basetypes.UsedFor_WithdrawalCompleted,
				Extra:     &extra,
				Vars:      notifVars,
				NotifType: basetypes.NotifType_NotifUnicast,
			}); err != nil {
				return err
			}

			notif, ok := notifMap[tx.ID]
			if !ok {
				continue
			}

			state := txnotifmgrpb.TxState_WaitNotified
			_, err = txnotifmwcli.UpdateTx(ctx, &txnotifmgrpb.TxReq{
				ID:         &notif.ID,
				NotifState: &state,
			})
			if err != nil {
				return err
			}
		}

		offset += limit
	}

	return nil
}

func waitNotified(ctx context.Context) error { // nolint
	offset := int32(0)
	limit := int32(1000)

	for {
		notifs, _, err := txnotifmwcli.GetTxs(ctx, &txnotifmgrpb.Conds{
			NotifState: &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(txnotifmgrpb.TxState_WaitNotified)},
			TxType:     &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(basetypes.TxType_TxWithdraw)},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(notifs) == 0 {
			break
		}

		for _, notif := range notifs {
			_notifs, _, err := notifmwcli.GetNotifs(ctx, &notifmwpb.Conds{
				Extra: &basetypes.StringVal{Op: cruder.LIKE, Value: notif.TxID},
			}, 0, int32(1000)) // nolint
			if err != nil {
				return err
			}
			if len(_notifs) == 0 {
				continue
			}

			notified := false
			for _, _notif := range _notifs {
				if !_notif.Notified {
					continue
				}
				notified = true
				break
			}

			if !notified {
				continue
			}

			state := txnotifmgrpb.TxState_Notified
			_, err = txnotifmwcli.UpdateTx(ctx, &txnotifmgrpb.TxReq{
				ID:         &notif.ID,
				NotifState: &state,
			})
			if err != nil {
				return err
			}
		}

		offset += limit
	}

	return nil
}

func processTx(ctx context.Context) {
	if err := waitSuccess(ctx); err != nil {
		logger.Sugar().Errorw("processTx", "error", err)
		return
	}
	if err := waitNotified(ctx); err != nil {
		logger.Sugar().Errorw("processTx", "error", err)
		return
	}
}
