package exec

import (
	"io"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/winebarrel/cwli/cli"
)

type Executable interface {
	Start(svc *cloudwatchlogs.CloudWatchLogs, flags *cli.Flags, out io.Writer) error
}

var Commands = []func(string) (Executable, error){
	parseHelp,
	parseShowGroups,
	parseHead,
	parseQuery,
}
