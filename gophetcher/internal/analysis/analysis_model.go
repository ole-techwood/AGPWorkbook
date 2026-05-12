package analysis

import (
	"strconv"

	"github.com/ole-techwood/AGPWorkbook/internal/report"
)

type AuditResult interface {
	ToReportRow() report.ReportRow
}

type BaseAuditResult struct {
	URL         string
	Status      string // HTTP status code as string or error marker.
	ErrorReason string
}

type WebAuditResult struct {
	BaseAuditResult

	ContentLength int64
	Server        string
}

type FileAuditResult struct {
	BaseAuditResult

	Size        int64
	Permissions string
}

type CriticalAuditResult struct {
	BaseAuditResult
}

func (r WebAuditResult) ToReportRow() report.ReportRow {
	size := "-"
	if _, err := strconv.Atoi(r.Status); err == nil {
		if r.ContentLength >= 0 {
			size = strconv.FormatInt(r.ContentLength, 10)
		}
	} else {
		size = report.SanitizeErrorMessage(r.ErrorReason)
	}

	return report.ReportRow{
		Resource: r.URL,
		Status:   r.Status,
		Size:     size,
		Info:     r.Server,
	}
}

func (r FileAuditResult) ToReportRow() report.ReportRow {
	size := report.SanitizeErrorMessage(r.ErrorReason)
	if r.Status == "OK" {
		size = strconv.FormatInt(r.Size, 10)
	}

	return report.ReportRow{
		Resource: r.URL,
		Status:   r.Status,
		Size:     size,
		Info:     r.Permissions,
	}
}

func (r CriticalAuditResult) ToReportRow() report.ReportRow {
	recoveredMessage := report.SanitizeErrorMessage(r.ErrorReason)
	if recoveredMessage == "" {
		recoveredMessage = "unknown panic"
	}

	return report.ReportRow{
		Resource: r.URL,
		Status:   "CRITICAL",
		Size:     "-",
		Info:     "CRITICAL FAILURE: RECOVERED: " + recoveredMessage,
	}
}
