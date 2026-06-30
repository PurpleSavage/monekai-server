package authports

type ResponseProvider struct{
	AuthID string
	Email string
	PhotoUrl *string
	PhoneNumber *string
}

type AuthProviderPort interface{
	VerifyAndExtract(token string) (*ResponseProvider,error)
}