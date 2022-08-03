package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	message struct {
		path   string
		format string
		args   []any
	}

	cmdOptions struct {
		dir        string
		name       string
		args       []string
		retries    int
		retryDelay int
	}
)

var (
	stdoutchan = make(chan message)
	done       = make(chan bool)
)

func main() {
	if len(os.Args) == 1 {
		showUsage()
		return
	}

	cmds := getCommandOptions(os.Args[1:])

	wg := sync.WaitGroup{}
	wg.Add(len(cmds))

	go func() {
		path := ""

		for {
			select {
			case <-done:
				return
			case msg := <-stdoutchan:
				if msg.path != path {
					path = msg.path
					fmt.Printf("\n======== %s ========\n", path)
				}

				fmt.Printf(msg.format, msg.args...)
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())

	for _, cmd := range cmds {
		go runCommand(ctx, &wg, cmd)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	go func() {
		<-sig
		cancel()
	}()

	wg.Wait()

	done <- true
}

func getCommandOptions(args []string) []cmdOptions {
	cmds := []cmdOptions{}
	dir := ""
	retries := 0
	retryDelay := 0

	for i, arg := range args {
		switch arg {
		case "--retries":
			retries, _ = strconv.Atoi(args[i+1])
		case "--retry-delay":
			retryDelay, _ = strconv.Atoi(args[i+1])
		case "--dir":
			dir = args[i+1]
		case "--cmd":
			tokens := strings.Split(args[i+1], " ")

			name := tokens[0]
			args := []string{}

			if len(tokens) > 1 {
				args = tokens[1:]
			}

			opts := cmdOptions{
				dir:        dir,
				name:       name,
				args:       args,
				retries:    retries,
				retryDelay: retryDelay,
			}

			cmds = append(cmds, opts)
		}
	}

	return cmds
}
func runCommand(ctx context.Context, wg *sync.WaitGroup, opts cmdOptions) {
	defer wg.Done()
	retries := 0

	for {
		cmd := exec.CommandContext(ctx, opts.name, opts.args...)
		cmd.Dir = opts.dir

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

		err = cmd.Wait()
		if err != nil {
			if err.Error()[:6] == "signal" {
				break
			}
		}

		stdoutchan <- message{
			path:   cmd.Path,
			format: "exited with code %d\n",
			args:   []any{cmd.ProcessState.ExitCode()},
		}

		if cmd.ProcessState.ExitCode() == 0 || opts.retries == 0 || opts.retries == retries {
			break
		}

		if opts.retryDelay > 0 {
			time.Sleep(time.Duration(opts.retryDelay) * time.Millisecond)
		}

		retries++
	}
}

func showUsage() {
	fmt.Println("Usage:")
	fmt.Println("  gather [--retries <retries>] [--retry-delay <delay>] [[--dir <dir>] --cmd <cmd>]...")
}
