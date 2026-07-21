package notificationsusecases

import (
	"context"
	notificationssreponsesdtos "github.com/PurpleSavage/monekai-server/modules/notifications/application/dtos/responses"
	notificationsports "github.com/PurpleSavage/monekai-server/modules/notifications/application/ports"
	commonentities "github.com/PurpleSavage/monekai-server/modules/shared/common/domain/entities"
	"golang.org/x/sync/errgroup"
)

type ListNotificationsUseCase struct {
	notificationsRepo notificationsports.NotificationsPersistencePort
}

func NewListNotificationsUseCase(
	notificationsRepo notificationsports.NotificationsPersistencePort,
) *ListNotificationsUseCase {
	return &ListNotificationsUseCase{
		notificationsRepo: notificationsRepo,
	}
}

func (uc *ListNotificationsUseCase) Execute(
	ctx context.Context,
	userID string,
	limit int,
	page int,
) (*commonentities.PaginatedResult[notificationssreponsesdtos.ItemNotificationDTO], error) {
	var totalNotifications int
	var notifications []notificationssreponsesdtos.ItemNotificationDTO
	g, ctx := errgroup.WithContext(ctx)
	
	g.Go(func() error {
		var err error
		totalNotifications, err = uc.notificationsRepo.CountTotalNotifications(ctx, userID)
		return err
	})
	
	g.Go(func() error {
		var err error
		notifications, err = uc.notificationsRepo.ListNotifications(ctx, userID, limit, page)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &commonentities.PaginatedResult[notificationssreponsesdtos.ItemNotificationDTO]{
		Total: totalNotifications,
		Limit: limit,
		Page:  page,
		Data:  notifications,
	}, nil
}
