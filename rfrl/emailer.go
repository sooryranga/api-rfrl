package rfrl

type VerificationEmailTypeAttributes struct {
	Heading     string
	Description string
}

var EmailTypeToEmailAttributes = map[string]VerificationEmailTypeAttributes{
	UserEmail: {
		Heading:     "Welcome To rfrl!",
		Description: "Please complete sign-up by using the passcode provided bellow:",
	},
	WorkEmail: {
		Heading:     "Verify Company to start referring!",
		Description: "Connect to company by using the passcode provided bellow:",
	},
}

type EmailerUseCase interface {
	SendEmailVerification(email string, emailType string) (string, error)
}
