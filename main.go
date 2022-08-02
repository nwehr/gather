package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type message struct {
	path   string
	format string
	args   []any
}

var (
	stdoutchan = make(chan message)
)

func main() {
	go func() {
		path := ""

		for msg := range stdoutchan {
			if msg.path != path {
				path = msg.path
				fmt.Printf("\n======== %s ========\n", path)
			}

			fmt.Printf(msg.format, msg.args...)
		}
	}()

	cmds := getCommands(os.Args[1:])

	wg := sync.WaitGroup{}

	for _, cmd := range cmds {
		wg.Add(1)
		go runCommand(&wg, cmd)
	}

	wg.Wait()
}

func getCommands(args []string) []*exec.Cmd {
	cmds := []*exec.Cmd{}
	dir := ""

	for i, arg := range args {
		switch arg {
		case "--dir":
			dir = args[i+1]
		case "--cmd":
			tokens := strings.Split(args[i+1], " ")

			name := tokens[0]
			args := []string{}

			if len(tokens) > 1 {
				args = tokens[1:]
			}

			cmd := exec.Command(name, args...)

			if dir != "" {
				cmd.Dir = dir
			}

			cmds = append(cmds, cmd)
		}
	}

	return cmds
}
func runCommand(wg *sync.WaitGroup, cmd *exec.Cmd) {
	defer wg.Done()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("could not get stdout pipe from %s: %v\n", cmd.String(), err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("could not get stdout pipe from %s: %v\n", cmd.String(), err)
	}

	err = cmd.Start()
	if err != nil {
		fmt.Printf("could not start %s: %v\n", cmd.String(), err)
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			stdoutchan <- message{
				path:   cmd.Path,
				format: "%s\n",
				args:   []any{scanner.Text()},
			}
		}
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		stdoutchan <- message{
			path:   cmd.Path,
			format: "%s\n",
			args:   []any{scanner.Text()},
		}
	}

	cmd.Wait()

	stdoutchan <- message{
		path:   cmd.Path,
		format: "exited with code %d\n",
		args:   []any{cmd.ProcessState.ExitCode()},
	}
}
