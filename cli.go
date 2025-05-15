package main

import (
	"fmt"
	"net/http"
	"net/url"
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

    _, err := url.ParseRequestURI(v)

    if err != nil {
      fmt.Printf("%s is not a valid url!\n", v)
      os.Exit(2)
    }

    filtered = append(filtered, v)
  }

  return filtered
}

func removeDuplicates(slice []string) []string {
    seen := make(map[string]struct{})
    result := []string{}

    for _, v := range slice {
        if _, found := seen[v]; !found {
            seen[v] = struct{}{}
            result = append(result, v)
        }
    }
    return result
}


func ParseArgs(args []string) []*Website {
  websiteUrlsDuplicate := checkArgs(args)
  websiteUrls := removeDuplicates(websiteUrlsDuplicate)

  websites := []*Website{}

  httpClient := &http.Client{
    Timeout: 10*time.Second,
  }

  for _, url := range websiteUrls {
    websites = append(websites, NewWebsite(url, httpClient))
  }

  return websites
}
