package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	var (
		trap  = flag.Bool("trap", false, "trap signals")
		sleep = flag.Duration("sleep", 0, "sleep time")
	)
	flag.Parse()

	if *trap {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		go func() {
			for _ = range sigCh {
				fmt.Fprintln(os.Stderr, "signal received")
			}
		}()
	}

	if *sleep > 0 {
		time.Sleep(*sleep)
	}

	fmt.Println(strings.Join(flag.Args(), " "))
}
