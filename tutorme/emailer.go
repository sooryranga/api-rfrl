package tutorme

type EmailerUseCase interface {
	SendEmailVerification(email string) (string, error)
}
