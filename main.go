package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	cmds := []*exec.Cmd{}
	dir := ""

	for i, arg := range os.Args {
		switch arg {
		case "--dir":
			dir = os.Args[i+1]
		case "--cmd":
			tokens := strings.Split(os.Args[i+1], " ")

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

	wg := sync.WaitGroup{}
	for i, cmd := range cmds {
		wg.Add(1)

		go func(i int, cmd *exec.Cmd) {
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
					fmt.Printf("%d: %s\n", i, scanner.Text())
				}
			}()

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Printf("%d: %s\n", i, scanner.Text())
			}
		}(i, cmd)
	}

	wg.Wait()
}
