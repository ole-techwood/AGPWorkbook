package report

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ReportHandler interface {
	Report(rows []ReportRow)
}

type ReportRow struct {
	Resource string `json:"resource"`
	Status   string `json:"status"`
	Size     string `json:"size"`
	Info     string `json:"info"`
}

type TextReportHandler struct{}

type JSONReportHandler struct{}

const TableRowFormat = "%-24s %-10s %-10s %s\n"

var ReportHandlers = map[string]ReportHandler{
	"text": TextReportHandler{},
	"json": JSONReportHandler{},
}

func NewReportHandler(format string) ReportHandler {
	requestedFormat := strings.ToLower(strings.TrimSpace(format))
	if handler, found := ReportHandlers[requestedFormat]; found {
		return handler
	}

	return ReportHandlers["text"]
}

func (rh TextReportHandler) Report(rows []ReportRow) {
	if len(rows) == 0 {
		fmt.Println("No audit results to display.")
		return
	}

	fmt.Println(strings.Repeat("-", 72))
	fmt.Printf(TableRowFormat, "RESOURCE", "STATUS", "SIZE", "INFO")
	fmt.Println(strings.Repeat("-", 72))

	for _, row := range rows {
		fmt.Printf(TableRowFormat, row.Resource, row.Status, row.Size, row.Info)
	}

	fmt.Println(strings.Repeat("-", 72))
	fmt.Printf("Done. %d resources processed.\n", len(rows))
}

func (rh JSONReportHandler) Report(rows []ReportRow) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(rows); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to write JSON report: %v\n", err)
	}
}

func SanitizeErrorMessage(message string) string {
	message = strings.ReplaceAll(message, "\n", " ")
	message = strings.ReplaceAll(message, "\r", " ")

	return strings.Join(strings.Fields(message), " ")
}
