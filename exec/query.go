package exec

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/pkg/errors"
	"github.com/winebarrel/cwli/cli"
)

const queryWaitInterval = 1

var regexpQuery = regexp.MustCompile(`(?i)^source\s+(\S+)\s+start=(\S+)\s+end=(\S+)\s*\|\s*(.+)$`)
var regexpQueryVertically = regexp.MustCompile(`(?i)\|\s*vertically\s*$`)

type queryCommand struct {
	logGroupName string
	queryString  string
	startTime    int64
	endTime      int64
	vertically   bool
}

func parseQuery(str string) (cmd Executable, err error) {
	submatch := regexpQuery.FindStringSubmatch(str)

	if len(submatch) == 0 {
		err = fmt.Errorf("Invalid query: %s", str)
		return
	}

	logGroupName := submatch[1]
	queryString := strings.TrimSpace(submatch[4])
	vertically := false

	if regexpQueryVertically.MatchString(queryString) {
		queryString = regexpQueryVertically.ReplaceAllString(queryString, "")
		queryString = strings.TrimSpace(queryString)
		vertically = true
	}

	var startTime, endTime time.Time
	startTime, err = dateparse.ParseLocal(submatch[2])

	if err != nil {
		err = errors.Wrap(err, "Could not parse startTime")
		return
	}

	endTime, err = dateparse.ParseLocal(submatch[3])

	if err != nil {
		err = errors.Wrap(err, "Could not parse endTime")
		return
	}

	cmd = &queryCommand{
		logGroupName: logGroupName,
		queryString:  queryString,
		startTime:    startTime.Unix(),
		endTime:      endTime.Unix(),
		vertically:   vertically,
	}

	return
}

func (cmd *queryCommand) startQuery(svc *cloudwatchlogs.CloudWatchLogs) (queryId *string, err error) {
	params := &cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(cmd.logGroupName),
		QueryString:  aws.String(cmd.queryString),
		StartTime:    aws.Int64(cmd.startTime),
		EndTime:      aws.Int64(cmd.endTime),
	}

	resp, err := svc.StartQuery(params)

	if err == nil {
		queryId = resp.QueryId
	}

	return
}

func waitQueryResult(svc *cloudwatchlogs.CloudWatchLogs, queryId *string) (result *cloudwatchlogs.GetQueryResultsOutput, err error) {
	params := &cloudwatchlogs.GetQueryResultsInput{
		QueryId: queryId,
	}

	for {
		result, err = svc.GetQueryResults(params)

		if err != nil {
			return
		}

		if *result.Status != "Scheduled" && *result.Status != "Running" {
			break
		}

		time.Sleep(queryWaitInterval * time.Second)
	}

	return
}

func escapeJson(str string) string {
	b, err := json.Marshal(str)

	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}

func printResultsVertically(results [][]*cloudwatchlogs.ResultField, showptr bool, out io.Writer) {
	if len(results) > 0 {
		for i, result := range results {
			fmt.Fprintf(out, "*************************** %d. row ***************************\n", i+1)

			for _, field := range result {
				if !showptr && *field.Field == "@ptr" {
					continue
				}

				fmt.Fprintf(out, "%s:\t%s\n", *field.Field, *field.Value)
			}
		}
	}
}

func printResultsHorizontally(results [][]*cloudwatchlogs.ResultField, showptr bool, out io.Writer) {
	if len(results) > 0 {
		fieldLen := len(results[0])

		if !showptr {
			fieldLen--
		}

		for _, result := range results {
			fmt.Fprint(out, "{")

			for i, field := range result {
				if !showptr && *field.Field == "@ptr" {
					continue
				}

				name := escapeJson(*field.Field)
				value := escapeJson(*field.Value)
				fmt.Fprintf(out, `%s:%s`, name, value)

				if i < fieldLen-1 {
					fmt.Fprint(out, ",")
				}
			}

			fmt.Fprintln(out, "}")
		}
	}
}

func (cmd *queryCommand) Start(svc *cloudwatchlogs.CloudWatchLogs, flags *cli.Flags, out io.Writer) (err error) {
	queryId, err := cmd.startQuery(svc)

	if err != nil {
		return
	}

	resp, err := waitQueryResult(svc, queryId)

	if err != nil {
		return
	}

	if cmd.vertically {
		printResultsVertically(resp.Results, flags.Showptr, out)
	} else {
		printResultsHorizontally(resp.Results, flags.Showptr, out)
	}

	fmt.Fprintf(out, "// Status: %s\n", *resp.Status)

	fmt.Fprintf(
		out,
		"// Statistics: BytesScanned=%.f RecordsMatched=%.f RecordsScanned=%.f\n",
		*resp.Statistics.BytesScanned,
		*resp.Statistics.RecordsMatched,
		*resp.Statistics.RecordsScanned,
	)

	return
}

func usageQuery() string {
	return "source LOG-GROUP start=START-TIME end=END-TIME | field ... [! EXTERNAL-COMMAND]\n\tPerform CloudWatch Logs Insights query"
}
