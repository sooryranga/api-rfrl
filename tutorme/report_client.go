package tutorme

import "time"

type ReportClient struct {
	ID        string    `db:"id" json:"id"`
	Reporter  string    `db:"reporter" json:"reporter"`
	Accused   string    `db:"accused" json:"accused"`
	Cause     string    `db:"cause" cause:"cause"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Tally     int       `db:"tally" json:"tally"`
}

func NewReportClient(reporter string, accused string, cause string) ReportClient {
	return ReportClient{
		Reporter: reporter,
		Accused:  accused,
		Cause:    cause,
	}
}

type ReportClientUseCase interface {
	CreateReport(report ReportClient) error
	DeleteReport(report ReportClient) error
	GetReports() (*[]ReportClient, error)
}

type ReportClientStore interface {
	CreateReport(db DB, report ReportClient) error
	DeleteReport(db DB, reporter string, accused string) error
	GetReports(db DB) (*[]ReportClient, error)
}
