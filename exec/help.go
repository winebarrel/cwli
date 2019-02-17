package exec

import (
	"fmt"
	"io"
	"regexp"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/winebarrel/cwli/cli"
)

var usages = []func() string{
	usageHelp,
	usageShowGroups,
	usageHead,
	usageQuery,
}

var regexpHelp = regexp.MustCompile(`(?i)^help\s*$`)

type helpCommand struct {
}

func parseHelp(str string) (cmd Executable, err error) {
	if !regexpHelp.MatchString(str) {
		return
	}

	cmd = &helpCommand{}

	return
}

func (cmd *helpCommand) Start(svc *cloudwatchlogs.CloudWatchLogs, flags *cli.Flags, out io.Writer) (err error) {

	for _, usage := range usages {
		fmt.Fprintln(out, usage())
	}

	return
}

func usageHelp() string {
	return "help\n\tPrint a help message"
}
