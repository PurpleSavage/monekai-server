package commonrequestsdtos
type ListNotificationsQueryDTO struct {
	Page  int `json:"page" validate:"required,min=1"`
	Limit int `json:"limit" validate:"required,min=1,max=100"`
}