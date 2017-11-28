package handlers

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
)

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing in directory: %s\n", cmd.Dir)
	fmt.Printf("==>   command line: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}