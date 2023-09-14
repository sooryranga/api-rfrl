package user

import "github.com/go-playground/validator"

type (
	// UserPayload is the struct used to hold payload from /user
	UserPayload struct {
		Email     string `json:"email" validate:"omitempty,email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Photo     string `json:"profileImageURL"`
		About     string `json:"about"`
	}

	// EducationPayload is the struct used to create education
	EducationPaylod struct {
		Institution     string `json:"institution"`
		Degree          string `json:"degree"`
		FeildOfStudy    string `json:"field_of_study"`
		start           string `json:"start" validate:"datetime"`
		end             string `json:"end" validate:"omitempty, datetime"`
		InstitutionLogo string `json:"institution_logo"`
	}
)

// SignUpPayloadValidation validates user inputs
func SignUpPayloadValidation(sl validator.StructLevel) {

	payload := sl.Current().Interface().(SignUpPayload)

	switch payload.Type {
	case GOOGLE:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case LINKEDIN:
		if len(payload.Token) == 0 {
			sl.ReportError(payload.Token, "token", "Token", "validToken", "")
		}
	case EMAIL:
		if len(payload.Email) == 0 {
			sl.ReportError(payload.Email, "email", "Email", "validEmail", "")
		}
		if len(payload.Password) < 10 {
			sl.ReportError(payload.Email, "password", "Password", "validPassworrd", "")
		}
	}

	// plus can do more, even with different tag than "fnameorlname"
}
