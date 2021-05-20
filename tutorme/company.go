package tutorme

import "gopkg.in/guregu/null.v4"

type Company struct {
	Name     string      `db:"company_name" json:"name"`
	Photo    null.String `db:"photo" json:"photo"`
	Industry null.String `db:"industry" json:"industry"`
	About    null.String `db:"about" json:"about"`
	Active   null.Bool   `db:"active" json:"-"`
}

func NewCompany(
	name string,
	photo null.String,
	industry null.String,
	about null.String,
	active null.Bool,
) Company {
	return Company{
		Name:     name,
		Photo:    photo,
		Industry: industry,
		About:    about,
		Active:   active,
	}
}

type CompanyUseCase interface {
	UpdateCompany(name string, photo null.String, industry null.String, about null.String, active null.Bool) (*Company, error)
	UpdateCompanyEmail(name string, emailDomain string, active bool) error
	GetCompanies(active bool) (*[]Company, error)
}

type CompanyStore interface {
	SelectCompany(db DB, name string) (*Company, error)
	CreateOrSelectCompany(db DB, name string) (*Company, error)
	CreateCompanyEmailDomain(db DB, emailDomain string) error
	UpdateCompany(db DB, company Company) (*Company, error)
	UpdateOrCreateCompanyEmail(db DB, name string, emailDomain string, active bool) error
	GetCompanies(db DB, active bool) (*[]Company, error)
	GetCompanyNameFromEmailDomain(db DB, domain string) (null.String, error)
}
