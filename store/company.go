package store

import (
	"errors"

	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"gopkg.in/guregu/null.v4"
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
INSERT INTO company_email (email_domain)
VALUES ($1)
`

func (cs *CompanyStore) CreateCompanyEmailDomain(db tutorme.DB, emailDomian string) error {
	_, err := db.Queryx(createCompanyEmailDomain, emailDomian)

	return err
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
INSERT INTO company_email (email_domain, company_id, active)
VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT email_domain 
DO UPDATE
SET company_name = $2, active = $3
`

func (cs *CompanyStore) UpdateOrCreateCompanyEmail(db tutorme.DB, name string, emailDomain string, active bool) error {
	company, err := cs.SelectCompany(db, name)

	if err != nil {
		return err
	}

	_, err = db.Queryx(updateOrCreateCompanyEmail, emailDomain, company.ID, active)
	return err
}

const getCompanies string = `
SELECT * FROM company WHERE active = $1
ORDER BY id
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

const getCompanyNameFromEmailDomainQuery string = `
SELECT company_name FROM company_email
INNER JOIN company ON company.id = company_email.company_id
WHERE email_domain = $1
`

func (cs *CompanyStore) GetCompanyNameFromEmailDomain(db tutorme.DB, emailDomain string) (null.String, error) {
	var name null.String

	err := db.QueryRowx(getCompanyNameFromEmailDomainQuery, emailDomain).Scan(&name)

	return name, err
}
