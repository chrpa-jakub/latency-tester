package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func checkArgs(args []string) []string {
  if(len(args) < 2) {
    panic(fmt.Sprintf("Usage: %s [url]", args[0]))
  }

  filtered := []string{}
  for _, v := range args[1:] {
      if v != "" && v != " " {

        if !strings.HasPrefix(v, "http://") && !strings.HasPrefix(v, "https://") {
          filtered = append(filtered, "http://"+v)
          continue
        }

      filtered = append(filtered, v)
    }

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
