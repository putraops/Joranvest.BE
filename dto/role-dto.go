package dto

type RoleAccessDto struct {
	RoleId    string `json:"role_id" binding:"required"`
	IsChecked bool   `json:"is_checked"`
}

type RoleNotificationDto struct {
	RoleId           string    `json:"role_id" binding:"required"`
	NotificationName *[]string `json:"notification_name"`
}

type PaymentNotificationDto struct {
	RoleId    string `json:"role_id" binding:"required"`
	IsChecked bool   `json:"is_checked"`
}
