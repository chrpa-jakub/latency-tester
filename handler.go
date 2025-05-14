package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartMeasuring(websites []*Website) {
  stop := make(chan os.Signal, 1)
  signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

  ticker := time.NewTicker(5 * time.Second)
  defer ticker.Stop()

  done := make(chan bool)

  for {
    select {
      case <-ticker.C: {
        go MeasureAllAsync(websites, done)
      }

      case <-stop: {
        return
      }

      case <-done: {
        go printAll(websites)
      }
    }
  }
}

func MeasureAllAsync(websites []*Website, done chan<-bool) {
  for _, website := range websites {
    go website.MeasureRequestAsync(done)
  } 
}

func printAll(websites []*Website) {
  fmt.Print("\033[H\033[2J")
  for _, website := range websites {
    website.Print()
  }
}
