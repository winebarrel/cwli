package cwli

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func executeCommand(str string, input string) (err error) {
	cmd := exec.Command("sh", "-c", str)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, input)
	}()

	err = cmd.Run()

	return
}
