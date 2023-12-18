package usecases

import (
	"github.com/Arun4rangan/api-tutorme/tutorme"
	"github.com/jmoiron/sqlx"
)

type ReportClientUseCase struct {
	db                *sqlx.DB
	ReportClientStore tutorme.ReportClientStore
}

func NewReportClientUseCase(db sqlx.DB, reportClientStore tutorme.ReportClientStore) *ReportClientUseCase {
	return &ReportClientUseCase{&db, reportClientStore}
}

func (r ReportClientUseCase) CreateReport(report tutorme.ReportClient) error {
	return r.ReportClientStore.CreateReport(r.db, report)
}

func (r ReportClientUseCase) DeleteReport(report tutorme.ReportClient) error {
	return r.ReportClientStore.DeleteReport(r.db, report.Reporter, report.Accused)
}

func (r ReportClientUseCase) GetReports() (*[]tutorme.ReportClient, error) {
	return r.ReportClientStore.GetReports(r.db)
}
