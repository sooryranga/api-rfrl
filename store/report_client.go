package store

import (
	tutorme "github.com/Arun4rangan/api-tutorme/tutorme"
)

// ReportClientStore holds all store related functions for reporting clients
type ReportClientStore struct{}

// NewReportClientStore creates new ReportClientStore
func NewReportClientStore() *ReportClientStore {
	return &ReportClientStore{}
}

const createReportQuery string = `
INSERT INTO report_client (reporter, accused, cause)
VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT report_client_reporter_accused_key 
DO UPDATE
SET tally = report_client.tally + 1
`

func (r ReportClientStore) CreateReport(db tutorme.DB, report tutorme.ReportClient) error {
	_, err := db.Queryx(createReportQuery, report.Reporter, report.Accused, report.Cause)
	return err
}

const deleteReportQuery string = `
DELETE FROM report_client
WHERE reporter = $1 AND accused = $2
`

func (r ReportClientStore) DeleteReport(db tutorme.DB, reporter string, accused string) error {
	_, err := db.Queryx(deleteReportQuery, reporter, accused)
	return err
}

const getReportQuery string = `
SELECT * FROM report_client
LIMIT 500
`

func (r ReportClientStore) GetReports(db tutorme.DB) (*[]tutorme.ReportClient, error) {
	rows, err := db.Queryx(getReportQuery)

	reports := make([]tutorme.ReportClient, 0)

	if err != nil {
		return &reports, err
	}

	for rows.Next() {
		var report tutorme.ReportClient

		err = rows.StructScan(&report)

		if err != nil {
			return &reports, err
		}
		reports = append(reports, report)
	}

	return &reports, nil
}
