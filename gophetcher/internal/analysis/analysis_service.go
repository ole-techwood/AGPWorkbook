package analysis

import (
	"context"
	"fmt"

	"github.com/ole-techwood/AGPWorkbook/internal/report"
)

// auditKey is a typed, unexported key for context values to avoid collisions.
type auditKey int

const auditRunIDKey auditKey = iota

// SetAuditRunID attaches an audit run ID to the context.
func SetAuditRunID(ctx context.Context, runID string) context.Context {
	return context.WithValue(ctx, auditRunIDKey, runID)
}

// GetAuditRunID retrieves the audit run ID from the context.
// Returns an empty string if no ID was set.
func GetAuditRunID(ctx context.Context) string {
	id, _ := ctx.Value(auditRunIDKey).(string)
	return id
}

type AnalysisService interface {
	AnalyzeURLs(ctx context.Context, urls []string) []AuditResult
}

func RecoverPanic(ctx context.Context, rawURL string, auditFn func() AuditResult) (result AuditResult) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = CriticalAuditResult{
				BaseAuditResult: BaseAuditResult{
					URL:         rawURL,
					Status:      "CRITICAL",
					ErrorReason: fmt.Sprint(recovered),
					RequestID:   GetAuditRunID(ctx),
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
