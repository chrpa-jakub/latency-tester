package main

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestWebsite_MeasureRequest(t *testing.T) {
  t.Parallel()

  httpmock.Activate()
  defer httpmock.DeactivateAndReset()

  tests := []struct {
    name          string
    mockResponder httpmock.Responder
    wantSuccess   uint64
    wantFail      uint64
    wantMinSize   int
    wantMaxSize   int
  }{
    {
      name: "Successful response 200 with small body",
      mockResponder: httpmock.NewStringResponder(200, "Hello"),
      wantSuccess: 1,
      wantFail:    0,
      wantMinSize: 5,
      wantMaxSize: 5,
    },
    {
      name: "Successful response 200 with larger body",
      mockResponder: httpmock.NewStringResponder(200, "Hello, this is a longer response body"),
      wantSuccess: 1,
      wantFail:    0,
      wantMinSize: 37,
      wantMaxSize: 37,
    },
    {
      name: "Failure due to 500 status",
      mockResponder: httpmock.NewStringResponder(500, "Internal Server Error"),
      wantSuccess: 0,
      wantFail:    1,
      wantMinSize: 0,
      wantMaxSize: 0,
    },
    {
      name: "Failure due to network error",
      mockResponder: httpmock.NewErrorResponder(assert.AnError),
      wantSuccess: 0,
      wantFail:    1,
      wantMinSize: 0,
      wantMaxSize: 0,
    },
  }

  for _, tt := range tests {
    tt := tt 
    t.Run(tt.name, func(t *testing.T) {
      t.Parallel() 

      url := "https://example.com/test"
      httpmock.RegisterResponder("GET", url, tt.mockResponder)

      w := NewWebsite(url, &http.Client{Transport: httpmock.DefaultTransport})

      w.MeasureRequest()

      assert.Equal(t, tt.wantSuccess, w.RequestInfo.SuccesCount)
      assert.Equal(t, tt.wantFail, w.RequestInfo.FailCount)

      if w.RequestInfo.SuccesCount > 0 {
        assert.Equal(t, tt.wantMinSize, w.SizeData.Min)
        assert.Equal(t, tt.wantMaxSize, w.SizeData.Max)
      }
    })
  }
}

