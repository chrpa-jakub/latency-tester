package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
	"time"
)

type Website struct {
	Url string
	LatencyData dataInfo
	SizeData dataInfo
	RequestInfo requestInfo
	httpClient *http.Client
	mu sync.Mutex
}

type requestInfo struct {
	SuccesCount uint64
	FailCount uint64
}

type dataInfo struct {
	Max int
	Avg float32
	Min int
}

func NewWebsite(url string, httpClient *http.Client) *Website {
	return &Website {
		Url: url,
		httpClient: httpClient,

		RequestInfo: requestInfo {
			SuccesCount: 0,
			FailCount: 0,
		},

		SizeData: dataInfo {
			Max: math.MinInt,
			Min: math.MaxInt,
			Avg: 0,
		},

		LatencyData: dataInfo {
			Max: math.MinInt,
			Min: math.MaxInt,
			Avg: 0,
		},
	}
}

func (r *requestInfo) RequestCount() uint64 {
	return r.SuccesCount + r.FailCount
}


func (w *Website) Print() {
	latency := w.LatencyData
	size := w.SizeData
	requests := w.RequestInfo

	if latency.Avg == 0 {
		fmt.Printf(
			"%s: Latency - Min: %d ms, Avg: %.2f ms, Max: %d ms | Size - Min: %d bytes, Avg: %.2f bytes, Max: %d bytes | Success: %d / %d\n",
			w.Url,
			0, latency.Avg, 0,
			0, latency.Avg, 0,
			requests.SuccesCount, requests.RequestCount(),
		)

		return
	}

	fmt.Printf(
		"%s: Latency - Min: %d ms, Avg: %.2f ms, Max: %d ms | Size - Min: %d bytes, Avg: %.2f bytes, Max: %d bytes | Success: %d / %d\n",
		w.Url,
		latency.Min, latency.Avg, latency.Max,
		size.Min, size.Avg, size.Max,
		requests.SuccesCount, requests.RequestCount(),
	)
}


func (w *Website) addMeasurements(size int, time int) {
	w.RequestInfo.SuccesCount++
	w.SizeData.addData(size, float32(w.RequestInfo.SuccesCount))
	w.LatencyData.addData(time, float32(w.RequestInfo.SuccesCount))
}

func (d *dataInfo) addData(data int, count float32) {
	d.Min = min(d.Min, data)
	d.Max = max(d.Max, data)
	d.Avg = ((count-1) * d.Avg + float32(data))/count
}

func (w *Website) MeasureRequest() *Website {
	w.mu.Lock()
	defer w.mu.Unlock()

	startTime := time.Now()
	resp, err := w.httpClient.Get(w.Url)

	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 400 {
		w.RequestInfo.FailCount++
		return w
	}

	defer resp.Body.Close()

	measuredTime := int(time.Since(startTime).Milliseconds())
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		w.RequestInfo.FailCount++
		return w
	}

	w.addMeasurements(len(body), measuredTime)
	return w
}
