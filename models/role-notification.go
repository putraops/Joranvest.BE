package models

type RoleNotification struct {
	BaseModel
	RoleId                  *string `gorm:"type:varchar(50);" json:"role_id"`
	HasPaymentNotification  *bool   `gorm:"type:bool;default:0" json:"has_payment_notification"`
	PaymentNotificationType *string `gorm:"type:varchar(20);" json:"payment_notification_type"`

	Role Role `gorm:"foreignkey:RoleId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
}

func (RoleNotification) TableName() string {
	return "role_notification"
}
