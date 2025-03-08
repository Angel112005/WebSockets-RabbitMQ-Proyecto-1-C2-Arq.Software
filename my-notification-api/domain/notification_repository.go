package domain

type NotificationRepository interface {
	Save(notification Notification) error
}
