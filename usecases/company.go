package usecases

import (
	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
)

type CompanyUseCase struct {
	db           *sqlx.DB
	CompanyStore tutorme.CompanyStore
}

func NewCompanyUseCase(db sqlx.DB, companyStore tutorme.CompanyStore) *CompanyUseCase {
	return &CompanyUseCase{&db, companyStore}
}

func (comu *CompanyUseCase) CreateSuggestion(name string, emailDomain string) (*tutorme.Company, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = comu.db.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	if *err != nil {
		return nil, *err
	}

	var company *tutorme.Company
	company, *err = comu.CompanyStore.CreateOrSelectCompany(tx, name)

	if *err != nil {
		return nil, *err
	}

	_, *err = comu.CompanyStore.CreateCompanyEmailDomain(tx, company.Name, emailDomain)

	if *err != nil {
		return nil, *err
	}

	return company, nil
}

func (comu *CompanyUseCase) UpdateCompany(
	name string,
	photo null.String,
	industry null.String,
	about null.String,
	active null.Bool,
) (*tutorme.Company, error) {
	var err = new(error)
	var tx *sqlx.Tx

	tx, *err = comu.db.Beginx()

	defer tutorme.HandleTransactions(tx, err)

	company := tutorme.NewCompany(name, photo, industry, about, active)
	var updatedCompany *tutorme.Company
	updatedCompany, *err = comu.CompanyStore.UpdateCompany(tx, company)

	if *err != nil {
		return nil, *err
	}

	return updatedCompany, nil
}

func (comu *CompanyUseCase) UpdateCompanyEmail(
	name string,
	emailDomain string,
	active bool,
) error {
	return comu.CompanyStore.UpdateOrCreateCompanyEmail(comu.db, name, emailDomain, active)
}

func (comu *CompanyUseCase) GetCompanies(active bool) (*[]tutorme.Company, error) {
	return comu.CompanyStore.GetCompanies(comu.db, active)
}
