package authresponsesdtos
type RenovateSessionResponseDto struct{
	UserData UserResponseDto `json:"userData"`
	AccessToken string `json:"accessToken"`
}