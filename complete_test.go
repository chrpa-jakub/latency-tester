package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

type errorReadCloser struct{}

func (e *errorReadCloser) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("read error")
}

func (e *errorReadCloser) Close() error {
	return nil
}

func TestMeasureRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		httpmock.Activate()

		url := "https://example.com"
		body := "Hello, world!"
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, body))

		client := &http.Client{Transport: httpmock.DefaultTransport}
		website := NewWebsite(url, client)

		result := website.MeasureRequest()

		if result.RequestInfo.SuccesCount != 1 {
			t.Errorf("Expected SuccessCount 1, got %d", result.RequestInfo.SuccesCount)
		}

		if result.RequestInfo.FailCount != 0 {
			t.Errorf("Expected FailCount 0, got %d", result.RequestInfo.FailCount)
		}

		if result.SizeData.Min != len(body) || result.SizeData.Max != len(body) {
			t.Errorf("Expected size min/max %d, got min=%d max=%d", len(body), result.SizeData.Min, result.SizeData.Max)
		}
	})

	t.Run("NetworkError", func(t *testing.T) {
		t.Parallel()

		httpmock.Activate()

		url := "https://fail.example.com"
		httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Failed")
		})

		client := &http.Client{Transport: httpmock.DefaultTransport}
		website := NewWebsite(url, client)

		result := website.MeasureRequest()

		if result.RequestInfo.SuccesCount != 0 {
			t.Errorf("Expected SuccessCount 0, got %d", result.RequestInfo.SuccesCount)
		}
		if result.RequestInfo.FailCount != 1 {
			t.Errorf("Expected FailCount 1, got %d", result.RequestInfo.FailCount)
		}
	})

	t.Run("HTTPErrorStatus", func(t *testing.T) {
		t.Parallel()

		httpmock.Activate()

		url := "https://error.example.com"
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(500, "Internal Server Error"))

		client := &http.Client{Transport: httpmock.DefaultTransport}
		website := NewWebsite(url, client)

		result := website.MeasureRequest()

		if result.RequestInfo.SuccesCount != 0 {
			t.Errorf("Expected SuccessCount 0, got %d", result.RequestInfo.SuccesCount)
		}
		if result.RequestInfo.FailCount != 1 {
			t.Errorf("Expected FailCount 1, got %d", result.RequestInfo.FailCount)
		}

	})

	t.Run("BodyReadError", func(t *testing.T) {
		t.Parallel()

		httpmock.Activate()

		url := "https://bodyerror.example.com"
		httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(202, "")
			res.Body = &errorReadCloser{}
			return res, nil
		})

		client := &http.Client{Transport: httpmock.DefaultTransport}
		website := NewWebsite(url, client)

		result := website.MeasureRequest()

		if result.RequestInfo.SuccesCount != 0 {
			t.Errorf("Expected SuccessCount 0, got %d", result.RequestInfo.SuccesCount)
		}
		if result.RequestInfo.FailCount != 1 {
			t.Errorf("Expected FailCount 1, got %d", result.RequestInfo.FailCount)
		}
	})

}

func TestMeasureRequest_ConcurrentMultipleURLs(t *testing.T) {
	t.Parallel()

	httpmock.Activate()

	urls := []string{
		"https://example1.com",
		"https://example2.com",
		"https://example3.com",
		"https://example4.com",
		"https://example5.com",
		"https://example6.com",
		"https://example7.com",
		"https://example8.com",
		"https://example9.com",
		"https://example10.com",
	}

	client := &http.Client{Transport: httpmock.DefaultTransport}

	for _, url := range urls {
		body := "Response from " + url
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, body))
	}


	for _, url := range urls {
		url := url 
		t.Run(url, func(t *testing.T) {
			t.Parallel()

			website := NewWebsite(url, client)
			website.MeasureRequest()

			if website.RequestInfo.SuccesCount != 1 {
				t.Errorf("Expected SuccessCount 1 for %s, got %d", url, website.RequestInfo.SuccesCount)
			}
			if website.RequestInfo.FailCount != 0 {
				t.Errorf("Expected FailCount 0 for %s, got %d", url, website.RequestInfo.FailCount)
			}
			if website.SizeData.Min <= 0 || website.SizeData.Max <= 0 {
				t.Errorf("Expected positive SizeData for %s, got min=%d max=%d", url, website.SizeData.Min, website.SizeData.Max)
			}
		})
	}
}

func TestStartMeasuring(t *testing.T) {
	t.Parallel()

	httpmock.Activate()
	url := "https://test-example-go-website.com"
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, "OK"))

	client := &http.Client{Transport: httpmock.DefaultTransport}
	website := NewWebsite(url, client)
	websites := []*Website{website}

	go StartMeasuring(websites)

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if website.RequestInfo.SuccesCount > 0 {
			return
		}
	}

	t.Errorf("Expected at least one successful measurement")
}

