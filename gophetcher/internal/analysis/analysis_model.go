package analysis

type AuditResult struct {
	URL           string
	ContentLength int64
	Server        string
	Status        string // "200", "ERROR", "UNREACHABLE", "TIMEOUT", etc.
	ErrorReason   string // error message if Status is not a valid HTTP code
}
