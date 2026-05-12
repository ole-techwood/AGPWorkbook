package analysis

import (
	"fmt"

	"github.com/ole-techwood/AGPWorkbook/internal/report"
)

type AnalysisService interface {
	AnalyzeURLs(urls []string) []AuditResult
}

func RecoverPanic(rawURL string, auditFn func() AuditResult) (result AuditResult) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = CriticalAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:         rawURL,
					Status:      "CRITICAL",
					ErrorReason: fmt.Sprint(recovered),
				},
			}
		}
	}()

	return auditFn()
}

func ToReportRows(results []AuditResult) []report.ReportRow {
	rows := make([]report.ReportRow, 0, len(results))
	for _, result := range results {
		rows = append(rows, result.ToReportRow())
	}

	return rows
}
