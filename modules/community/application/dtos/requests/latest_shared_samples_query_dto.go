package communityrequestsdtos

type LatestSharedSamplesQueryDTO struct {
	Limit int `json:"limit" validate:"required,min=1,max=15"`
}
