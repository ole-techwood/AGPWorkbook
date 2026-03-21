package analyzer

type AuditResult struct {
	URL           string
	StatusCode    int
	ContentLength int64
	Server        string
}
