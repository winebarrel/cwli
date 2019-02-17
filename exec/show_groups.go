package exec

import (
	"fmt"
	"io"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/winebarrel/cwli/cli"
)

var regexpShowGroups = regexp.MustCompile(`(?i)^show\s+groups(?:\s+prefix\s+(\S+))?\s*$`)

type showGroupsCommand struct {
	prefix string
}

func parseShowGroups(str string) (query Executable, err error) {
	submatch := regexpShowGroups.FindStringSubmatch(str)

	if len(submatch) == 0 {
		return
	}

	query = &showGroupsCommand{
		prefix: submatch[1],
	}

	return
}

func (cmd *showGroupsCommand) Start(svc *cloudwatchlogs.CloudWatchLogs, flags *cli.Flags, out io.Writer) (err error) {
	params := &cloudwatchlogs.DescribeLogGroupsInput{}

	if cmd.prefix != "" {
		params.LogGroupNamePrefix = aws.String(cmd.prefix)
	}

	err = svc.DescribeLogGroupsPages(params, func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
		for _, group := range page.LogGroups {
			fmt.Fprintln(out, *group.LogGroupName)
		}

		return !lastPage
	})

	return
}

func usageShowGroups() string {
	return "show groups [prefix LOG-GROUP-PREFIX]\n\tDisplay Log Groups"
}
