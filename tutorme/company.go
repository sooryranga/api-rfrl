package tutorme

import "gopkg.in/guregu/null.v4"

type Company struct {
	ID       int         `db:"id" json:"id"`
	Name     null.String `db:"company_name" json:"name"`
	Photo    null.String `db:"photo" json:"photo"`
	Industry null.String `db:"industry" json:"industry"`
	About    null.String `db:"about" json:"about"`
	Active   null.Bool   `db:"active" json:"active"`
}

type CompanyEmailDomain struct {
	EmailDomain string   `db:"email_domain" json:"emailDomain"`
	CompanyID   null.Int `db:"company_id" json:"companyId"`
	Active      bool     `db:"active" json:"active"`
}

func NewCompany(
	name null.String,
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
	UpdateCompany(id int, name null.String, photo null.String, industry null.String, about null.String, active null.Bool) (*Company, error)
	CreateCompany(name string, photo null.String, industry null.String, about null.String, active null.Bool) (*Company, error)
	UpdateCompanyEmail(name string, emailDomain string, active bool) error
	GetCompanies(active bool) (*[]Company, error)
	GetCompany(id int) (*Company, error)
	GetCompanyEmails(withCompany null.Bool) (*[]CompanyEmailDomain, error)
}

type CompanyStore interface {
	SelectCompany(db DB, name string) (*Company, error)
	CreateOrSelectCompany(db DB, name string) (*Company, error)
	CreateCompanyEmailDomain(db DB, emailDomain string) error
	CreateCompany(db DB, company Company) (*Company, error)
	UpdateCompany(db DB, company Company) (*Company, error)
	UpdateOrCreateCompanyEmail(db DB, name string, emailDomain string, active bool) error
	GetCompanies(db DB, active bool) (*[]Company, error)
	GetCompanyIDFromEmailDomain(db DB, domain string) (null.Int, error)
	GetCompany(db DB, id int) (*Company, error)
	GetCompanyEmails(db DB, withCompany null.Bool) (*[]CompanyEmailDomain, error)
}
