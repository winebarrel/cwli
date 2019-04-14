package exec

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/winebarrel/cwli/cli"
)

var regexpHead = regexp.MustCompile(`(?i)^head\s+(\S+)(?:\s+limit\s+([0-9]+))?\s*$`)

type headCommand struct {
	group string
	limit int64
}

func parseHead(str string) (cmd Executable, err error) {
	submatch := regexpHead.FindStringSubmatch(str)

	if len(submatch) == 0 {
		return
	}

	group := submatch[1]
	limit := int64(3)

	if submatch[2] != "" {
		limit, err = strconv.ParseInt(submatch[2], 10, 64)

		if err != nil {
			return
		}
	}

	if limit <= 0 {
		err = fmt.Errorf("Invalid limit: %d", limit)
		return
	}

	cmd = &headCommand{
		group: group,
		limit: limit,
	}

	return
}

func describeLogGroup(svc *cloudwatchlogs.CloudWatchLogs, group string) (logGroup *cloudwatchlogs.LogGroup, err error) {
	params := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(group),
		Limit:              aws.Int64(1),
	}

	resp, err := svc.DescribeLogGroups(params)

	if err != nil {
		return
	}

	if len(resp.LogGroups) == 0 || *resp.LogGroups[0].LogGroupName != group {
		err = fmt.Errorf("LogGroup was not found: %s", group)
		return
	}

	logGroup = resp.LogGroups[0]

	return
}

func describeLogStream(svc *cloudwatchlogs.CloudWatchLogs, logGroup *cloudwatchlogs.LogGroup) (logStream *cloudwatchlogs.LogStream, err error) {
	params := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: logGroup.LogGroupName,
		Descending:   aws.Bool(true),
		OrderBy:      aws.String("LastEventTime"),
		Limit:        aws.Int64(1),
	}

	resp, err := svc.DescribeLogStreams(params)

	if err != nil {
		return
	}

	if len(resp.LogStreams) == 0 {
		err = fmt.Errorf("LogStream was not found in %s", *logGroup.LogGroupName)
		return
	}

	logStream = resp.LogStreams[0]

	return
}

func getLogEvent(svc *cloudwatchlogs.CloudWatchLogs, logGroup *cloudwatchlogs.LogGroup, logStream *cloudwatchlogs.LogStream, limit int64) (events []*cloudwatchlogs.OutputLogEvent, err error) {
	params := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  logGroup.LogGroupName,
		LogStreamName: logStream.LogStreamName,
		StartTime:     logStream.LastEventTimestamp,
		Limit:         aws.Int64(limit),
	}

	resp, err := svc.GetLogEvents(params)

	if err != nil {
		return
	}

	if len(resp.Events) == 0 {
		err = fmt.Errorf("Events was not found in %s", *logGroup.LogGroupName)
		return
	}

	events = resp.Events

	return
}

func (cmd *headCommand) Start(svc *cloudwatchlogs.CloudWatchLogs, flags *cli.Flags, out io.Writer) (err error) {
	logGroup, err := describeLogGroup(svc, cmd.group)

	if err != nil {
		return
	}

	logStream, err := describeLogStream(svc, logGroup)

	if err != nil {
		return
	}

	events, err := getLogEvent(svc, logGroup, logStream, cmd.limit)

	if err != nil {
		return
	}

	for i, event := range events {
		fmt.Fprintf(out, "*************************** %d. row ***************************\n", i+1)
		fmt.Fprintf(out, "Timestamp:\t%d\n", *event.Timestamp)
		fmt.Fprintf(out, "Message:\t%s\n", *event.Message)
	}

	return
}

func usageHead() string {
	return "head LOG-GROUP [limit N]\n\tPrint first lines of a Log Group"
}
