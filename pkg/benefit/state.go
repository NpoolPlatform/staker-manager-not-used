package benefit

import (
	gbmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/goodbenefit"
	pltfaccmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/platform"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"
	goodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"

	"github.com/shopspring/decimal"
)

type State struct {
	Coins            map[string]*coinmwpb.Coin
	PlatformAccounts map[string]map[basetypes.AccountUsedFor]*pltfaccmwpb.Account // map[CoinTypeID]map[UsedFor]Account
	GoodBenefits     map[string]*gbmwpb.Account
	ChangeState      bool
	UpdateGoodProfit bool
}

func newState() *State {
	return &State{
		Coins:            map[string]*coinmwpb.Coin{},
		PlatformAccounts: map[string]map[basetypes.AccountUsedFor]*pltfaccmwpb.Account{},
		GoodBenefits:     map[string]*gbmwpb.Account{},
		ChangeState:      true,
		UpdateGoodProfit: true,
	}
}

type Good struct {
	*goodmwpb.Good
	TodayRewardAmount         decimal.Decimal
	PlatformRewardAmount      decimal.Decimal
	UserRewardAmount          decimal.Decimal
	TechniqueServiceFeeAmount decimal.Decimal
	BenefitOrderIDs           map[string]string
	BenefitOrderUnits         decimal.Decimal
	BenefitAccountAmount      decimal.Decimal
	Retry                     bool
}

func newGood(good *goodmwpb.Good) *Good {
	return &Good{
		good,
		decimal.NewFromInt(0),
		decimal.NewFromInt(0),
		decimal.NewFromInt(0),
		decimal.NewFromInt(0),
		map[string]string{},
		decimal.NewFromInt(0),
		decimal.NewFromInt(0),
		false,
	}
}
