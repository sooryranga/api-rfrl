package usecases

import (
	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type CompanyUseCase struct {
	db           *sqlx.DB
	CompanyStore rfrl.CompanyStore
}

func NewCompanyUseCase(db sqlx.DB, companyStore rfrl.CompanyStore) *CompanyUseCase {
	return &CompanyUseCase{&db, companyStore}
}

func (comu *CompanyUseCase) GetCompany(id int) (*rfrl.Company, error) {
	return comu.CompanyStore.GetCompany(comu.db, id)
}

func (comu *CompanyUseCase) GetCompanyEmails(withCompany null.Bool) (*[]rfrl.CompanyEmailDomain, error) {
	return comu.CompanyStore.GetCompanyEmails(comu.db, withCompany)
}

func (comu *CompanyUseCase) CreateCompany(
	name string,
	photo null.String,
	industry null.String,
	about null.String,
	active null.Bool,
) (*rfrl.Company, error) {

	company := rfrl.NewCompany(null.StringFrom(name), photo, industry, about, active)
	return comu.CompanyStore.CreateCompany(comu.db, company)
}

func (comu *CompanyUseCase) UpdateCompany(
	ID int,
	name null.String,
	photo null.String,
	industry null.String,
	about null.String,
	active null.Bool,
) (*rfrl.Company, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = comu.db.Beginx()

	defer rfrl.HandleTransactions(tx, err)

	company := rfrl.NewCompany(name, photo, industry, about, active)
	company.ID = ID
	var updatedCompany *rfrl.Company
	updatedCompany, *err = comu.CompanyStore.UpdateCompany(tx, company)

	if *err != nil {
		return nil, *err
	}

	return updatedCompany, nil
}

func (comu *CompanyUseCase) UpdateCompanyEmail(
	name null.String,
	emailDomain string,
	active bool,
) error {
	return comu.CompanyStore.UpdateOrCreateCompanyEmail(comu.db, name, emailDomain, active)
}

func (comu *CompanyUseCase) GetCompanies(active bool) (*[]rfrl.Company, error) {
	return comu.CompanyStore.GetCompanies(comu.db, active)
}

func (comu *CompanyUseCase) GetCompanyEmail(companyEmail string) (*rfrl.CompanyEmailDomain, error) {
	return comu.CompanyStore.GetCompanyEmail(comu.db, companyEmail)
}
