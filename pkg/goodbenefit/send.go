package goodbenefit

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	appgoodcli "github.com/NpoolPlatform/good-middleware/pkg/client/appgood"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	"github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
	appgoodpb "github.com/NpoolPlatform/message/npool/good/mw/v1/appgood"
	notifmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif"
	notifbenefitpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/notif/goodbenefit"
	tmplmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/template"
	notifmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif"
	notifbenefitcli "github.com/NpoolPlatform/notif-middleware/pkg/client/notif/goodbenefit"
	"github.com/shopspring/decimal"
)

//nolint
func send(ctx context.Context, channel basetypes.NotifChannel) {
	offset := int32(0)
	limit := int32(1000)

	for {
		goodBenefits, _, err := notifbenefitcli.GetGoodBenefits(ctx, &notifbenefitpb.Conds{
			Generated: &basetypes.BoolVal{Op: cruder.EQ, Value: false},
		}, offset, limit)
		if err != nil {
			logger.Sugar().Errorw("GetGoodBenefits:", "Error", err)
			return
		}
		if len(goodBenefits) == 0 {
			logger.Sugar().Info("goodbenefits:", "length:", 0)
			break
		}

		benefitIDs := []string{}
		goodIDs := []string{}
		for _, benefit := range goodBenefits {
			goodIDs = append(goodIDs, benefit.GoodID)
			benefitIDs = append(benefitIDs, benefit.ID)
		}

		goods, _, err := appgoodcli.GetGoods(ctx, &appgood.Conds{
			GoodIDs: &npool.StringSliceVal{
				Op:    cruder.IN,
				Value: goodIDs,
			},
		}, 0, 10000)
		if err != nil {
			logger.Sugar().Errorw("GetGoods", "error", err)
		}
		if len(goods) == 0 {
			logger.Sugar().Errorw("GetGoods", "length", len(goods))
			break
		}

		goodMap := map[string]*appgoodpb.Good{}
		for _, _good := range goods {
			goodMap[_good.GoodID] = _good
		}

		content := "<html><head><style>table.notif-benefit {border-collapse: collapse;width: 100%;}#notif-good-benefit td,#notif-good-benefit th {border: 1px solid #dddddd;text-align: left;padding: 8px;}</style></head><table id='notif-good-benefit' class='notif-benefit'><tr><th>GoodID</th><th>GoodName</th><th>Amount</th><th>AmountPerUnit</th><th>State</th><th>Message</th><th>BenefitDate</th></tr>" //nolint

		for _, benefit := range goodBenefits {
			good, ok := goodMap[benefit.GoodID]
			if !ok {
				logger.Sugar().Errorw("goodBenefits", "goodMap", ok)
				continue
			}

			total, err := decimal.NewFromString(good.GoodTotal)
			if err != nil {
				logger.Sugar().Errorw("GoodTotal", "err", err)
			}
			amount, err := decimal.NewFromString(benefit.Amount)
			if err != nil {
				logger.Sugar().Errorw("Amount", "err", err)
			}

			tm := time.Unix(int64(benefit.BenefitDate), 0)
			content += fmt.Sprintf(`<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>`,
				benefit.GoodID, benefit.GoodName,
				benefit.Amount, amount.Div(total), benefit.State,
				benefit.Message, tm,
			)
		}
		content += "</table></html>"

		logger.Sugar().Infow("Content", content)

		// find AppID from Goods
		appIDs := map[string]struct{}{}
		for _, _good := range goods {
			appIDs[_good.AppID] = struct{}{}
		}

		if err != nil {
			logger.Sugar().Errorf("Marshal", "Error", err)
		}

		for appID, _ := range appIDs {
			_, err := notifmwcli.GenerateNotifs(ctx, &notifmwpb.GenerateNotifsRequest{
				AppID:     appID,
				EventType: basetypes.UsedFor_GoodBenefit1,
				NotifType: basetypes.NotifType_NotifMulticast,
				Vars: &tmplmwpb.TemplateVars{
					Message: &content,
				},
			})
			if err != nil {
				logger.Sugar().Errorw("GenerateNotifs", "Error", err)
			}

			generated := true
			for _, benefitID := range benefitIDs {
				_, err := notifbenefitcli.UpdateGoodBenefit(ctx, &notifbenefitpb.GoodBenefitReq{
					ID:        &benefitID,
					Generated: &generated,
				})
				if err != nil {
					logger.Sugar().Errorw("UpdateGoodBenefit", "Error", err)
				}
			}
		}

		offset += limit
	}
}
