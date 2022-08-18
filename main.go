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

var (
	commit string
	built  string
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
		wait       int

		ctx context.Context
	}
)

func main() {
	if len(os.Args) == 1 {
		showUsage()
		return
	}

	output := make(chan message)
	done := make(chan bool)

	go printOutput(output, done)

	ctx, _ := getContextWithCancel()
	cmdOptions := getCommandOptions(ctx, os.Args[1:])

	wg := sync.WaitGroup{}
	wg.Add(len(cmdOptions))

	for _, opts := range cmdOptions {
		go runCommand(&wg, output, opts)
	}

	wg.Wait()
	done <- true
}

func getContextWithCancel() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := make(chan os.Signal, 4)
		signal.Notify(sig, os.Interrupt)
		<-sig
		cancel()
	}()

	return ctx, cancel
}

func getCommandOptions(ctx context.Context, args []string) []cmdOptions {
	cmds := []cmdOptions{}

	dir := ""
	retries := 0
	retryDelay := 0
	wait := 0

	for i, arg := range args {
		switch arg {
		case "--retries":
			retries, _ = strconv.Atoi(args[i+1])
		case "--retry-delay":
			retryDelay, _ = strconv.Atoi(args[i+1])
		case "--wait":
			wait, _ = strconv.Atoi(args[i+1])

		case "-d":
			fallthrough
		case "--dir":
			dir = args[i+1]

		case "-c":
			fallthrough
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
				wait:       wait,
				ctx:        ctx,
			}

			cmds = append(cmds, opts)
		}
	}

	return cmds
}

func runCommand(wg *sync.WaitGroup, output chan message, opts cmdOptions) {
	defer wg.Done()

	for {
		cmd := exec.CommandContext(opts.ctx, opts.name, opts.args...)
		cmd.Dir = opts.dir

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("could not get stdout pipe from %s: %v\n", cmd.String(), err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Printf("could not get stdout pipe from %s: %v\n", cmd.String(), err)
		}

		if opts.wait > 0 {
			time.Sleep(time.Duration(opts.wait) * time.Millisecond)
		}

		err = cmd.Start()
		if err != nil {
			fmt.Printf("could not start %s: %v\n", cmd.String(), err)
		}

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				output <- message{
					path:   cmd.Path,
					format: "%s\n",
					args:   []any{scanner.Text()},
				}
			}
		}()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			output <- message{
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

		output <- message{
			path:   cmd.Path,
			format: "exited with code %d\n",
			args:   []any{cmd.ProcessState.ExitCode()},
		}

		if cmd.ProcessState.ExitCode() == 0 || opts.retries == 0 {
			break
		}

		if opts.retryDelay > 0 {
			time.Sleep(time.Duration(opts.retryDelay) * time.Millisecond)
		}

		opts.retries--
	}
}

func printOutput(output chan message, done chan bool) {
	path := ""

	for {
		select {
		case <-done:
			return
		case msg := <-output:
			if msg.path != path {
				path = msg.path
				fmt.Printf("\n=======> %s <=======\n\n", path)
			}

			fmt.Printf(msg.format, msg.args...)
		}
	}
}

func showUsage() {
	fmt.Println("gather:")
	fmt.Printf("  version: %s\n", commit)
	fmt.Printf("  build:   %s\n", built)
	fmt.Println()

	fmt.Println("Options:")
	fmt.Println("  --retries <retries>    Optional: Number of times to retry failed cmd")
	fmt.Println("  --retry-delay <delay>  Optional: Wait time in ms before each retry")
	fmt.Println("  --wait <delay>         Optional: Wait time in ms before cmd is started")
	fmt.Println("  --dir | -d <dir>       Optional: Working dir of cmd")
	fmt.Println("  --cmd | -c <cmd>")
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Println("  gather [options...]")
	fmt.Println()

	fmt.Println("Examples:")
	fmt.Println("  gather --cmd 'start_db.sh' --cmd 'start_server.sh'")
}
