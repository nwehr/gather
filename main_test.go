package main

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGetComandOptions(t *testing.T) {
	args := []string{"--cmd", "ls", "--dir", "./", "--retries", "3", "--retry-delay", "500", "--cmd", "mv gather /bin/"}
	cmdOptions := getCommandOptions(context.Background(), args)

	if len(cmdOptions) != 2 {
		t.Errorf("expected 2 commands, got %d", len(cmdOptions))
	}

	ls := cmdOptions[0]
	mv := cmdOptions[1]

	{
		if ls.retries != 0 {
			t.Errorf("expected 0 retries, got %d", ls.retries)
		}

		if ls.retryDelay != 0 {
			t.Errorf("expected 0 retry delay, got %d", ls.retryDelay)
		}

		if ls.dir != "" {
			t.Errorf("expected dir to be '', got '%s'", ls.dir)
		}

		if len(ls.args) != 0 {
			t.Errorf("expected 0 args, got %d", len(ls.args))
		}
	}

	{
		if mv.retries != 3 {
			t.Errorf("expected 3 retries, got %d", mv.retries)
		}

		if mv.retryDelay != 500 {
			t.Errorf("expected 500 retry delay, got %d", mv.retryDelay)
		}

		if mv.dir != "./" {
			t.Errorf("expected dir to be './', got '%s'", mv.dir)
		}

		if len(mv.args) != 2 {
			t.Errorf("expected 2 args, got %d", len(mv.args))
		}
	}
}

func TestRunCommand(t *testing.T) {
	{
		opts := cmdOptions{
			name: "echo",
			args: []string{"hello"},
			ctx:  context.Background(),
		}

		output := make(chan message)

		go runCommand(&sync.WaitGroup{}, output, opts)

		msg := <-output

		if strings.TrimSpace(fmt.Sprintf(msg.format, msg.args...)) != "hello" {
			t.Errorf("expected 'hello', got '%s'", fmt.Sprintf(msg.format, msg.args...))
		}
	}

	{
		opts := cmdOptions{
			name: "sleep",
			args: []string{"3"},
			ctx:  context.Background(),
		}

		output := make(chan message)

		wg := sync.WaitGroup{}
		wg.Add(1)

		start := time.Now()

		go runCommand(&wg, output, opts)

		<-output

		wg.Wait()

		stop := time.Now()

		duration := stop.Sub(start).Seconds()

		if math.Abs(duration-3.0) > 0.05 {
			t.Errorf("expected sleep to run for 3 seconds, got %f", duration)
		}
	}

	{
		ctx, cancel := context.WithCancel(context.Background())

		opts := cmdOptions{
			name: "sleep",
			args: []string{"5"},
			ctx:  ctx,
		}

		output := make(chan message)

		wg := sync.WaitGroup{}
		wg.Add(1)

		start := time.Now()

		go func(cancel context.CancelFunc) {
			time.Sleep(time.Second)
			cancel()
		}(cancel)

		go runCommand(&wg, output, opts)

		wg.Wait()

		stop := time.Now()

		duration := stop.Sub(start).Seconds()

		if math.Abs(duration-1.0) > 0.05 {
			t.Errorf("expected cancel after 1 second, got %f", duration)
		}
	}
}
