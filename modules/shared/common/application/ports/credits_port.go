package commonports

type CreditsPort interface {
	CheckCredits(email string, creditsNeeded int) (bool,error)
	DecreaseCredits(credits int,email string) (int,error)
}
