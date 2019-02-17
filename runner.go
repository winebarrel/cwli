package cwli

import (
	"bytes"
	"fmt"
	"log"
	"os/user"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/chzyer/readline"
	"github.com/winebarrel/cwli/cli"
)

// Runner struct has information on CloudWatch and Command-Line Flags.
type Runner struct {
	svc   *cloudwatchlogs.CloudWatchLogs
	flags *cli.Flags
}

// NewRunner creates Runner struct.
func NewRunner(flags *cli.Flags) (runner *Runner) {
	sess := session.Must(session.NewSession())

	runner = &Runner{
		svc:   cloudwatchlogs.New(sess),
		flags: flags,
	}

	return
}

// Run CloudWatch Logs Insights client.
func (runner *Runner) Run() {
	prompt := runner.svc.Client.ClientInfo.SigningRegion + "> "

	runQuery := func(str string) {
		var out bytes.Buffer
		query, cmd, err := parseCommand(str)

		if err == nil {
			err = query.Start(runner.svc, runner.flags, &out)
		}

		if err != nil {
			fmt.Printf("error: %s\n", err)
		}

		if cmd == "" {
			fmt.Print(out.String())
		} else {
			cmdErr := executeCommand(cmd, out.String())

			if cmdErr != nil {
				fmt.Printf("error: %s\n", cmdErr)
			}
		}
	}

	if runner.flags.Query != "" {
		query := runner.flags.Query

		if !strings.HasSuffix(query, ";") {
			query += ";"
		}

		runQuery(query)
	} else {
		repl(prompt, runQuery)
	}
}

func repl(prompt string, block func(str string)) {
	var buf strings.Builder

	currUser, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:      prompt,
		HistoryFile: currUser.HomeDir + "/.cwli_history",
	})

	if err != nil {
		panic(err)
	}

	defer rl.Close()

	for {
		line, err := rl.Readline()

		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		stmts := strings.SplitAfter(line, ";")

		for _, stmt := range stmts {
			if stmt == "" {
				continue
			}

			buf.WriteString(stmt)

			if !strings.HasSuffix(stmt, ";") {
				buf.WriteString(" ")
				continue
			}

			block(buf.String())
			buf.Reset()
		}
	}
}
