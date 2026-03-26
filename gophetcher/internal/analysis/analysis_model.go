package analysis

type AuditResult any

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
