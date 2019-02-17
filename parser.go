package cwli

import (
	"strings"

	"github.com/winebarrel/cwli/exec"
)

func parseCommand(str string) (query exec.Executable, cmd string, err error) {
	str = strings.TrimRight(str, ";")

	if n := strings.Index(str, "!"); n >= 0 {
		cmd = strings.TrimSpace(str[n+1:])
		str = str[0:n]
	}

	for _, fn := range exec.Commands {
		query, err = fn(str)

		if err != nil {
			return
		}

		if query != nil {
			return
		}
	}

	return
}
