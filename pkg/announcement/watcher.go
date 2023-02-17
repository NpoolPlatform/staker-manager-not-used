package announcement

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	usermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuser"
	usermwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	applangmgrpb "github.com/NpoolPlatform/message/npool/g11n/mgr/v1/applang"
	ancmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement"
	ancusermgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement/user"
	chanmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"
	emailtmplmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/email"
	smstmplmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/template/sms"
	ancmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	ancsendmwpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/sendstate"
	sendmwpb "github.com/NpoolPlatform/message/npool/third/mw/v1/send"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	applangmwcli "github.com/NpoolPlatform/g11n-middleware/pkg/client/applang"
	ancmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"
	ancsendmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/sendstate"
	ancusermwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"
	emailtmplmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/email"
	smstmplmwcli "github.com/NpoolPlatform/notif-middleware/pkg/client/template/sms"
	sendmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/send"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func unicast(ctx context.Context, anc *ancmwpb.Announcement, user *usermwpb.User) (bool, error) {
	req := &sendmwpb.SendMessageRequest{
		Subject: anc.Title,
		Content: anc.Content,
	}

	lang, err := applangmwcli.GetLangOnly(ctx, &applangmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: anc.AppID,
		},
		Main: &commonpb.BoolVal{
			Op:    cruder.EQ,
			Value: true,
		},
	})
	if err != nil {
		return false, err
	}
	if lang == nil {
		return false, fmt.Errorf("applang main invalid")
	}

	if lang.LangID != anc.LangID {
		return false, nil
	}

	switch anc.Channel {
	case chanmgrpb.NotifChannel_ChannelEmail:
		tmpl, err := emailtmplmwcli.GetEmailTemplateOnly(ctx, &emailtmplmgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: anc.AppID,
			},
			LangID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: anc.LangID,
			},
			UsedFor: &commonpb.Int32Val{
				Op:    cruder.EQ,
				Value: int32(basetypes.UsedFor_Announcement),
			},
		})
		if err != nil {
			return false, err
		}
		if tmpl == nil {
			return false, fmt.Errorf("email template invalid")
		}

		req.From = tmpl.Sender
		req.To = user.EmailAddress
		req.ToCCs = tmpl.CCTos
		req.ReplyTos = tmpl.ReplyTos
		req.AccountType = basetypes.SignMethod_Email
	case chanmgrpb.NotifChannel_ChannelSMS:
		tmpl, err := smstmplmwcli.GetSMSTemplateOnly(ctx, &smstmplmgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: anc.AppID,
			},
			LangID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: anc.LangID,
			},
			UsedFor: &commonpb.Int32Val{
				Op:    cruder.EQ,
				Value: int32(basetypes.UsedFor_Announcement),
			},
		})
		if err != nil {
			return false, err
		}
		if tmpl == nil {
			return false, fmt.Errorf("sms template invalid")
		}

		req.To = user.PhoneNO
		req.AccountType = basetypes.SignMethod_Mobile
	}

	if err := sendmwcli.SendMessage(ctx, req); err != nil {
		return false, err
	}
	return true, nil
}

func multicastUsers(ctx context.Context, anc *ancmwpb.Announcement, users []*usermwpb.User) error {
	uids := []string{}
	for _, user := range users {
		uids = append(uids, user.ID)
	}

	stats, _, err := ancsendmwcli.GetSendStates(ctx, &ancsendmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: anc.AppID,
		},
		AnnouncementID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: anc.AnnouncementID,
		},
		Channel: &commonpb.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(anc.Channel.Number()),
		},
		UserIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: uids,
		},
	}, 0, int32(len(uids)))
	if err != nil {
		return err
	}

	statMap := map[string]*ancsendmwpb.SendState{}
	for _, stat := range stats {
		statMap[stat.UserID] = stat
	}

	for _, user := range users {
		if _, ok := statMap[user.ID]; ok {
			logger.Sugar().Infow(
				"multicastUsers",
				"AppID", user.AppID,
				"UserID", user.ID,
				"EmailAddress", user.EmailAddress,
				"AnnouncementID", anc.AnnouncementID,
				"AnnoucementType", anc.AnnouncementType,
				"State", "Sent")
			continue
		}

		switch anc.Channel {
		case chanmgrpb.NotifChannel_ChannelEmail:
			if !strings.Contains(user.EmailAddress, "@") {
				logger.Sugar().Errorw(
					"multicastUsers",
					"AppID", user.AppID,
					"UserID", user.ID,
					"EmailAddress", user.EmailAddress,
					"State", "Invalid")
				continue
			}
		case chanmgrpb.NotifChannel_ChannelSMS:
			if user.PhoneNO == "" {
				logger.Sugar().Errorw(
					"multicastUsers",
					"AppID", user.AppID,
					"UserID", user.ID,
					"PhoneNO", user.PhoneNO,
					"State", "Invalid")
				continue
			}
		default:
			continue
		}

		_, err := unicast(ctx, anc, user)
		if err != nil {
			logger.Sugar().Errorw(
				"multicastUsers",
				"AppID", user.AppID,
				"UserID", user.ID,
				"EmailAddress", user.EmailAddress,
				"PhoneNO", user.PhoneNO,
				"error", err)
			return err
		}

		// TODO: record send state
	}

	return nil
}

func broadcast(ctx context.Context, anc *ancmwpb.Announcement) error {
	offset := int32(0)
	limit := int32(1000)

	for {
		users, _, err := usermwcli.GetUsers(ctx, &usermgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: anc.AppID,
			},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			break
		}

		if err := multicastUsers(ctx, anc, users); err != nil {
			return err
		}

		offset += limit
	}

	return nil
}

func multicast(ctx context.Context, anc *ancmwpb.Announcement) error {
	offset := int32(0)
	limit := int32(1000)

	for {
		ancUsers, _, err := ancusermwcli.GetUsers(ctx, &ancusermgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: anc.AppID,
			},
			AnnouncementID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: anc.AnnouncementID,
			},
		}, offset, limit)
		if err != nil {
			return err
		}
		if len(ancUsers) == 0 {
			break
		}

		uids := []string{}
		for _, user := range ancUsers {
			uids = append(uids, user.UserID)
		}

		users, _, err := usermwcli.GetManyUsers(ctx, uids)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			continue
		}

		if err := multicastUsers(ctx, anc, users); err != nil {
			return err
		}

		offset += limit
	}

	return nil
}

func sendOne(ctx context.Context, anc *ancmwpb.Announcement) error {
	switch anc.AnnouncementType {
	case ancmgrpb.AnnouncementType_Broadcast:
		return broadcast(ctx, anc)
	case ancmgrpb.AnnouncementType_Multicast:
		return multicast(ctx, anc)
	}
	return fmt.Errorf("announcement invalid")
}

func send(ctx context.Context, channel chanmgrpb.NotifChannel) {
	offset := int32(0)
	limit := int32(100)
	now := uint32(time.Now().Unix())

	for {
		ancs, _, err := ancmwcli.GetAnnouncements(ctx, &ancmwpb.Conds{
			EndAt: &commonpb.Uint32Val{
				Op:    cruder.GT,
				Value: now,
			},
			Channel: &commonpb.Uint32Val{
				Op:    cruder.EQ,
				Value: uint32(channel),
			},
		}, offset, limit)
		if err != nil {
			logger.Sugar().Errorw("send", "error", err)
			return
		}
		if len(ancs) == 0 {
			break
		}

		for _, anc := range ancs {
			if err := sendOne(ctx, anc); err != nil {
				logger.Sugar().Errorw("send", "error", err)
			}
		}

		offset += limit
	}
}

func Watch(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			send(ctx, chanmgrpb.NotifChannel_ChannelEmail)
			send(ctx, chanmgrpb.NotifChannel_ChannelSMS)
		case <-ctx.Done():
			return
		}
	}
}
