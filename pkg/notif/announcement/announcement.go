package announcement

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"

	notifmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/notif"

	usercli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v2/appuser"
	userpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/user"
	channelpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/channel"

	announcemenmgrtpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement"
	announcementpb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement"
	announcementcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement"

	userannouncementmgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement/user"
	userannouncementcli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/user"

	sendstatemgrpb "github.com/NpoolPlatform/message/npool/notif/mgr/v1/announcement/sendstate"
	sendstatepb "github.com/NpoolPlatform/message/npool/notif/mw/v1/announcement/sendstate"
	sendstatecli "github.com/NpoolPlatform/notif-middleware/pkg/client/announcement/sendstate"

	thirdpb "github.com/NpoolPlatform/message/npool/third/mgr/v1/template/notif"
	thirdcli "github.com/NpoolPlatform/third-middleware/pkg/client/notif"
	thirdtempcli "github.com/NpoolPlatform/third-middleware/pkg/client/template/notif"
)

var channel = channelpb.NotifChannel_ChannelEmail

//nolint:gocognit
func sendAnnouncement(ctx context.Context) {
	offset := int32(0)
	limit := int32(1000)
	endAt := uint32(time.Now().Unix())

	for {
		aInfos, _, err := announcementcli.GetAnnouncements(ctx, &announcemenmgrtpb.Conds{
			EndAt: &commonpb.Uint32Val{
				Op:    cruder.GT,
				Value: endAt,
			},
			Channels: &commonpb.StringSliceVal{
				Op:    cruder.IN,
				Value: []string{channel.String()},
			},
		}, offset, limit)
		if err != nil {
			logger.Sugar().Errorw("sendAnnouncement", "offset", offset, "limit", limit, "error", err)
			return
		}

		fmt.Println("***********aInfos")
		fmt.Println(aInfos)

		offset += limit
		if len(aInfos) == 0 {
			logger.Sugar().Info("There are no announcements to send within the announcement end at")
			return
		}

		for _, val := range aInfos {
			switch val.AnnouncementType {
			case announcemenmgrtpb.AnnouncementType_AllUsers:
				allUserUsersType(ctx, val)
			case announcemenmgrtpb.AnnouncementType_AppointUsers:
				appointUsersType(ctx, val)
			}
		}
	}
}

func appointUsersType(ctx context.Context, info *announcementpb.Announcement) {
	uOffset := int32(0)
	uLimit := int32(50)
	for {
		uAInfos, _, err := userannouncementcli.GetUsers(ctx, &userannouncementmgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: info.AppID,
			},
			AnnouncementID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: info.AnnouncementID,
			},
		}, uOffset, uLimit)
		if err != nil {
			logger.Sugar().Errorw("sendAnnouncement", "offset", uOffset, "limit", uLimit, "error", err)
			return
		}

		if len(uAInfos) == 0 {
			break
		}

		userID := []string{}
		for _, val := range uAInfos {
			userID = append(userID, val.UserID)
		}

		userInfos, _, err := usercli.GetManyUsers(ctx, userID)
		if err != nil {
			logger.Sugar().Errorw("sendAnnouncement", "offset", uOffset, "limit", uLimit, "error", err)
			return
		}
		uOffset += uLimit
		if len(userInfos) == 0 {
			continue
		}

		sendEmail(ctx, info, userInfos)
	}
}

func sendEmail(
	ctx context.Context,
	info *announcementpb.Announcement,
	userInfos []*userpb.User,
) {
	userMap := map[string]*userpb.User{}
	userIDs := []string{}
	for _, user := range userInfos {
		userIDs = append(userIDs, user.ID)
		userMap[user.ID] = user
	}

	templateInfo, err := thirdtempcli.GetNotifTemplateOnly(ctx, &thirdpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: info.AppID,
		},
		UsedFor: &commonpb.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(notifmgrpb.EventType_Announcement),
		},
	})
	if err != nil {
		logger.Sugar().Errorw("sendAnnouncement", "error", err)
		return
	}

	if templateInfo == nil {
		logger.Sugar().Errorw("sendAnnouncement", "AppID", info.AppID, "error", "template is empty")
		return
	}

	sendAnnou, _, err := sendstatecli.GetSendStates(ctx, &sendstatepb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: info.AppID,
		},
		AnnouncementID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: info.AnnouncementID,
		},
		Channel: &commonpb.Uint32Val{
			Op:    cruder.EQ,
			Value: uint32(channel.Number()),
		},
		UserIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: userIDs,
		},
	}, 0, int32(len(userIDs)))
	if err != nil {
		logger.Sugar().Errorw("sendAnnouncement", "error", err)
		return
	}

	sendAnnouMap := map[string]*announcementpb.Announcement{}
	for _, send := range sendAnnou {
		sendAnnouMap[send.UserID] = info
	}

	sendInfos := []*sendstatemgrpb.SendStateReq{}

	for _, user := range userInfos {
		if !strings.Contains(user.EmailAddress, "@") {
			continue
		}

		if _, ok := sendAnnouMap[user.ID]; ok {
			logger.Sugar().Infow(
				"sendAnnouncement",
				"AppID", user.AppID,
				"UserID", user.ID,
				"EmailAddress", user.EmailAddress,
				"AnnouncementID", info.AnnouncementID,
				"State", "Sent")
			continue
		}

		logger.Sugar().Infow("sendAnnouncement", "EmailAddress", user.EmailAddress, "AnnouncementID", info.AnnouncementID, "State", "Sending")
		err = thirdcli.SendNotifEmail(ctx, info.Title, info.Content, templateInfo.Sender, user.EmailAddress)
		if err != nil {
			logger.Sugar().Errorw("sendAnnouncement", "error", err.Error(), "Sender", templateInfo.Sender, "To", user.EmailAddress)
			return
		}
		sendInfos = append(sendInfos, &sendstatemgrpb.SendStateReq{
			AppID:          &info.AppID,
			UserID:         &user.ID,
			AnnouncementID: &info.AnnouncementID,
			Channel:        &channel,
		})
	}

	if len(sendInfos) == 0 {
		return
	}
	err = sendstatecli.CreateSendStates(ctx, sendInfos)
	if err != nil {
		logger.Sugar().Errorw("sendAnnouncement", "error", err.Error())
		return
	}
}

func allUserUsersType(ctx context.Context, info *announcementpb.Announcement) {
	uOffset := int32(0)
	uLimit := int32(50)
	for {
		userInfos, _, err := usercli.GetUsers(ctx, &appusermgrpb.Conds{
			AppID: &commonpb.StringVal{
				Op:    cruder.EQ,
				Value: info.AppID,
			},
		}, uOffset, uLimit)
		if err != nil {
			logger.Sugar().Errorw("sendAnnouncement", "offset", uOffset, "limit", uLimit, "error", err)
			return
		}
		uOffset += uLimit
		if len(userInfos) == 0 {
			break
		}
		sendEmail(ctx, info, userInfos)
	}
}
