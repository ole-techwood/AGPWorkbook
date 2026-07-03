package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/ole-techwood/AGPWorkbook/gorge/internal/game"
	"github.com/ole-techwood/AGPWorkbook/gorge/internal/pipeline"
)

// ExecuteStream runs hardcoded game analytics pipeline for provided stream.
func ExecuteStream(stream string) Report {
	rawEvents := strings.Split(stream, ",")
	report := Report{TotalRawEvents: len(rawEvents), Outputs: make([]string, 0, len(rawEvents))}

	parse := pipeline.Stage[string, game.GameEvent](game.ParseEvent)
	filter := pipeline.Filter(func(e game.GameEvent) bool { return e.Score > 0 })
	transform := pipeline.Map(game.ToSummary)
	pipe := pipeline.New(pipeline.Link(pipeline.Link(parse, filter), transform))

	for _, raw := range rawEvents {
		out, err := pipe.Run(raw)

		if err != nil {
			if errors.Is(err, pipeline.ErrFiltered) {
				report.FilteredOut++
			} else {
				report.ParseFailed++
			}
			continue
		}

		report.Outputs = append(report.Outputs, out)
	}

	return report
}

// Run executes `gorge run` command.
func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	stream := fs.String("stream", "", "comma-separated raw events")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, "failed to parse run flags:", err)
		return 1
	}

	if *stream == "" {
		fmt.Fprintln(stderr, "missing required flag: --stream")
		return 1
	}

	report := ExecuteStream(*stream)
	fmt.Fprint(stdout, report.String())
	return 0
}
