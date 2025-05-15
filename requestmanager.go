package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/inancgumus/screen"
)

func StartMeasuring(websites []*Website) {
  stop := make(chan os.Signal, 1)
  signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

  ticker := time.NewTicker(5 * time.Second)

  var wg sync.WaitGroup

  done := make(chan *Website)

  for {
    select {
      case <-ticker.C: {
        go func() {
          MeasureAllAsync(websites, done, &wg)
        }()
      }

      case <-stop: {
        go func() {
          ticker.Stop()
          wg.Wait()
          printAll(websites)
          os.Exit(0)
        }()
      }

      case <-done: {
        printAll(websites)
      }
    }
  }
}

func MeasureAllAsync(websites []*Website, done chan<-*Website, wg *sync.WaitGroup) {
  for _, website := range websites {
    wg.Add(1)
    time.Sleep(time.Millisecond*100)
    go func() {
      defer wg.Done()
      done<-website.MeasureRequest()
    }()
  } 
}

func printAll(websites []*Website) {
  screen.Clear()
  screen.MoveTopLeft()
  for _, website := range websites {
    website.Print()
  }
}
