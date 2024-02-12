package store

import (
	"errors"

	rfrl "github.com/Arun4rangan/api-rfrl/rfrl"
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

func (cs *CompanyStore) SelectCompany(db rfrl.DB, name string) (*rfrl.Company, error) {
	var company rfrl.Company

	err := db.QueryRowx(selectCompany, name).StructScan(&company)

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

func (cs *CompanyStore) CreateOrSelectCompany(db rfrl.DB, name string) (*rfrl.Company, error) {
	var company rfrl.Company

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

func (cs *CompanyStore) CreateCompanyEmailDomain(db rfrl.DB, emailDomian string) error {
	_, err := db.Queryx(createCompanyEmailDomain, emailDomian)

	return err
}

const createCompanyQuery string = `
INSERT INTO company (company_name, photo, industry, about, active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *
`

func (cs *CompanyStore) CreateCompany(db rfrl.DB, company rfrl.Company) (*rfrl.Company, error) {
	row := db.QueryRowx(
		createCompanyQuery,
		company.Name,
		company.Photo,
		company.Industry,
		company.About,
		company.Active,
	)
	var createdCompany rfrl.Company
	err := row.StructScan(&createdCompany)

	return &createdCompany, err
}

func (cs *CompanyStore) UpdateCompany(db rfrl.DB, company rfrl.Company) (*rfrl.Company, error) {
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

	if company.Name.Valid {
		query = query.Set("company_name", company.Name)
	}

	sql, args, err := query.
		Where(sq.Eq{"id": company.ID}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	var updatedCompany rfrl.Company
	err = db.QueryRowx(sql, args...).StructScan(&updatedCompany)

	return &updatedCompany, err
}

const updateOrCreateCompanyEmail string = `
INSERT INTO company_email (email_domain, company_id, active)
VALUES ($1, $2, $3)
ON CONFLICT (email_domain)
DO UPDATE
SET company_id = $2, active = $3
`

func (cs *CompanyStore) UpdateOrCreateCompanyEmail(db rfrl.DB, name null.String, emailDomain string, active bool) error {
	var companyId null.Int
	if name.Valid {
		company, err := cs.SelectCompany(db, name.String)

		if err != nil {
			return err
		}
		companyId = null.IntFrom(int64(company.ID))
	}

	_, err := db.Queryx(updateOrCreateCompanyEmail, emailDomain, companyId, active)
	return err
}

const getCompanies string = `
SELECT * FROM company WHERE active = $1
ORDER BY id
`

func (cs *CompanyStore) GetCompanies(db rfrl.DB, active bool) (*[]rfrl.Company, error) {
	companies := make([]rfrl.Company, 0)

	rows, err := db.Queryx(getCompanies, active)

	if err != nil {
		return &companies, err
	}

	for rows.Next() {
		var company rfrl.Company
		err := rows.StructScan(&company)
		if err != nil {
			return &companies, err
		}
		companies = append(companies, company)
	}
	return &companies, nil
}

const getCompanyIDFromEmailDomainQuery string = `
SELECT company_id FROM company_email
WHERE email_domain = $1
`

func (cs *CompanyStore) GetCompanyIDFromEmailDomain(db rfrl.DB, emailDomain string) (null.Int, error) {
	var id null.Int

	err := db.QueryRowx(getCompanyIDFromEmailDomainQuery, emailDomain).Scan(&id)

	return id, err
}

const getCompanyQuery string = `
SELECT * FROM company
WHERE id = $1
`

func (cs *CompanyStore) GetCompany(db rfrl.DB, id int) (*rfrl.Company, error) {
	var company rfrl.Company

	row := db.QueryRowx(getCompanyQuery, id)
	err := row.StructScan(&company)

	return &company, err
}

func (cs *CompanyStore) GetCompanyEmails(db rfrl.DB, withCompany null.Bool) (*[]rfrl.CompanyEmailDomain, error) {
	companyEmails := make([]rfrl.CompanyEmailDomain, 0)

	query := sq.Select("*").From("company_email")

	if withCompany.Valid && withCompany.Bool {
		query = query.Where(sq.NotEq{"company_id": nil})
	}

	if withCompany.Valid && !withCompany.Bool {
		query = query.Where(sq.Eq{"company_id": nil})
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return &companyEmails, err
	}

	rows, err := db.Queryx(sql, args...)

	if err != nil {
		return &companyEmails, err
	}

	for rows.Next() {
		var companyEmail rfrl.CompanyEmailDomain

		err = rows.StructScan(&companyEmail)

		if err != nil {
			return &companyEmails, err
		}

		companyEmails = append(companyEmails, companyEmail)
	}

	return &companyEmails, nil
}

const getCompanyEmailQuery string = `
SELECT * FROM company_email
WHERE email_domain = $1
`

func (cs *CompanyStore) GetCompanyEmail(db rfrl.DB, emailDomain string) (*rfrl.CompanyEmailDomain, error) {
	row := db.QueryRowx(getCompanyEmailQuery, emailDomain)

	var companyEmail rfrl.CompanyEmailDomain

	err := row.StructScan(&companyEmail)

	return &companyEmail, err
}
