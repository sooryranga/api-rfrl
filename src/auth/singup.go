package auth

type (
	// SignUpPayload is the struct used to hold payload from /signup
	SignUpPayload struct {
		Email           string `json:"email" validate:"email"`
		Name            string `json:"name"`
		Token           string `json:"token"`
		ProfileImageURL string `json:"profileImageURL"`
		Password        string `json:"password" validate:"gte=10,required_with=email"`
		Type            string `json:"type" validate:"required,oneof= GOOGLE LINKEDIN EMAIL"`
	}
)
