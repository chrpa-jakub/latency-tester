package main

import (
  "fmt"
  "net/http"
  "os"
  "strings"
  "time"
)

func checkArgs(args []string) []string {
  if(len(args) < 2) {
    fmt.Fprintf(os.Stderr, "Usage: %s [url]\n", args[0])
    os.Exit(1)
  }

  filtered := []string{}
  for _, v := range args[1:] {
    v = strings.TrimSpace(v)

    if v == "" {
      continue
    }

    if !strings.HasPrefix(v, "http://") && !strings.HasPrefix(v, "https://") {
    v = "http://"+v
  }

  filtered = append(filtered, v)
}

return filtered
}

func ParseArgs(args []string) []*Website {
  websiteUrls := checkArgs(args)

  websites := []*Website{}

  httpClient := &http.Client{
    Timeout: 10*time.Second,
  }

  for _, url := range websiteUrls {
    websites = append(websites, NewWebsite(url, httpClient))
  }

  return websites
}
