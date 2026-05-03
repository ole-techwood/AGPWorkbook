package analysis

import (
	"strconv"

	"github.com/ole-techwood/AGPWorkbook/internal/report"
)

type AnalysisService interface {
	AnalyzeURLs(urls []string) []AuditResult
}

func ToReportRows(results []AuditResult) []report.ReportRow {
	rows := make([]report.ReportRow, 0, len(results))
	for _, result := range results {
		rows = append(rows, toReportRow(result))
	}

	return rows
}

func toReportRow(result AuditResult) report.ReportRow {
	switch typed := result.(type) {
	case WebAuditResult:
		return webReportRow(typed)
	case *WebAuditResult:
		if typed == nil {
			return report.ReportRow{Resource: "-", Status: "ERROR", Size: "-", Info: "unexpected nil web result"}
		}
		return webReportRow(*typed)
	case FileAuditResult:
		return fileReportRow(typed)
	case *FileAuditResult:
		if typed == nil {
			return report.ReportRow{Resource: "-", Status: "ERROR", Size: "-", Info: "unexpected nil file result"}
		}
		return fileReportRow(*typed)
	default:
		return report.ReportRow{Resource: "-", Status: "ERROR", Size: "-", Info: "unexpected result type"}
	}
}

func webReportRow(result WebAuditResult) report.ReportRow {
	size := "-"
	if _, err := strconv.Atoi(result.Status); err == nil {
		if result.ContentLength >= 0 {
			size = strconv.FormatInt(result.ContentLength, 10)
		}
	} else {
		size = report.SanitizeErrorMessage(result.ErrorReason)
	}

	return report.ReportRow{
		Resource: result.URL,
		Status:   result.Status,
		Size:     size,
		Info:     result.Server,
	}
}

func fileReportRow(result FileAuditResult) report.ReportRow {
	size := report.SanitizeErrorMessage(result.ErrorReason)
	if result.Status == "OK" {
		size = strconv.FormatInt(result.Size, 10)
	}

	return report.ReportRow{
		Resource: result.URL,
		Status:   result.Status,
		Size:     size,
		Info:     result.Permissions,
	}
}
