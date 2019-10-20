package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
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
		sigs := []os.Signal{os.Interrupt}
		if runtime.GOOS != "windows" {
			sigs = append(sigs, syscall.SIGTERM)
		}
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, sigs...)
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
