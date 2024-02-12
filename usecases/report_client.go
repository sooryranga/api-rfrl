package usecases

import (
	"github.com/Arun4rangan/api-rfrl/rfrl"
	"github.com/jmoiron/sqlx"
)

type ReportClientUseCase struct {
	db                *sqlx.DB
	ReportClientStore rfrl.ReportClientStore
}

func NewReportClientUseCase(db sqlx.DB, reportClientStore rfrl.ReportClientStore) *ReportClientUseCase {
	return &ReportClientUseCase{&db, reportClientStore}
}

func (r ReportClientUseCase) CreateReport(report rfrl.ReportClient) error {
	return r.ReportClientStore.CreateReport(r.db, report)
}

func (r ReportClientUseCase) DeleteReport(report rfrl.ReportClient) error {
	return r.ReportClientStore.DeleteReport(r.db, report.Reporter, report.Accused)
}

func (r ReportClientUseCase) GetReports() (*[]rfrl.ReportClient, error) {
	return r.ReportClientStore.GetReports(r.db)
}
