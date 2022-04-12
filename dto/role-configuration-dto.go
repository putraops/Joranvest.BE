package dto

type RoleNotificationDto struct {
	RoleId           string    `json:"role_id" binding:"required"`
	NotificationName *[]string `json:"notification_name"`
}
