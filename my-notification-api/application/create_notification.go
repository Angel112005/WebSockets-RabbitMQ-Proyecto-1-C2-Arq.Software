package application

import "my-notification-api/domain"

type CreateNotification struct {
	repo domain.NotificationRepository
}

func NewCreateNotification(repo domain.NotificationRepository) *CreateNotification {
	return &CreateNotification{repo: repo}
}

func (uc *CreateNotification) Execute(notification domain.Notification) error {
	return uc.repo.Save(notification)
}
