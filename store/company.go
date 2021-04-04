package store

import (
	"errors"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// CompanyStore holds all store related functions for company
type CompanyStore struct{}

// NewCompanyStore creates new clientStore
func NewCompanyStore() *CompanyStore {
	return &CompanyStore{}
}

const selectCompany string = `
SELECT * FROM company WHERE company_name = $1 
`

func (cs *CompanyStore) SelectCompany(db tutorme.DB, name string) (*tutorme.Company, error) {
	var company tutorme.Company

	err := db.QueryRowx(createCompany, name).StructScan(&company)

	if err != nil {
		return nil, err
	}

	return &company, nil
}

const createCompany string = `
INSERT INTO company (company_name)
VALUES ($1)
RETURNING *
`

func (cs *CompanyStore) CreateOrSelectCompany(db tutorme.DB, name string) (*tutorme.Company, error) {
	var company tutorme.Company

	err := db.QueryRowx(createCompany, name).StructScan(&company)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return cs.SelectCompany(db, name)
			}
		}
		return nil, err
	}

	return &company, err
}

const createCompanyEmailDomain string = `
INSERT INTO company_email (company_name, email_domain)
VALUES ($1, $2)
ON CONFLICT (email_domain) DO UPDATE
SET suggestions = company_email.suggestions + 1;
`

func (cs *CompanyStore) CreateCompanyEmailDomain(db tutorme.DB, name string, emailDomian string) (string, error) {
	_, err := db.Queryx(createCompanyEmailDomain, name, emailDomian)

	if err != nil {
		return "", err
	}

	return emailDomian, nil
}

func (cs *CompanyStore) UpdateCompany(db tutorme.DB, company tutorme.Company) (*tutorme.Company, error) {
	query := sq.Update("company")

	if company.Photo.Valid {
		query = query.Set("photo", company.Photo)
	}

	if company.Industry.Valid {
		query = query.Set("industry", company.Industry)
	}

	if company.About.Valid {
		query = query.Set("about", company.About)
	}

	if company.Active.Valid {
		query = query.Set("active", company.Active)
	}

	sql, args, err := query.Suffix("RETURNING *").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	var updatedCompany tutorme.Company
	err = db.QueryRowx(sql, args...).StructScan(&updatedCompany)

	return &updatedCompany, err
}

const updateOrCreateCompanyEmail string = `
INSERT INTO company_email (email_domain, company_name, active)
VALUES ($1, $2, $3)
ON CONFLICT (email_domain) DO UPDATE
SET company_name = $2, active = $3
`

func (cs *CompanyStore) UpdateOrCreateCompanyEmail(db tutorme.DB, name string, emailDomain string, active bool) error {
	_, err := db.Queryx(updateOrCreateCompanyEmail, emailDomain, name, active)
	return err
}

const getCompanies string = `
SELECT * FROM company WHERE active = $1
ORDER BY company_name
`

func (cs *CompanyStore) GetCompanies(db tutorme.DB, active bool) (*[]tutorme.Company, error) {
	companies := make([]tutorme.Company, 0)

	rows, err := db.Queryx(getCompanies, active)

	if err != nil {
		return &companies, err
	}

	for rows.Next() {
		var company tutorme.Company
		err := rows.StructScan(&company)
		if err != nil {
			return &companies, err
		}
		companies = append(companies, company)
	}
	return &companies, nil
}
