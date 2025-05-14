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

	if(w.LatencyData.Avg == 0) {
		fmt.Printf("%s: Latency - Avg: %.2f, Max: %d, Min: %d Size - Avg: %.2f, Max: %d, Min: %d %d/%d \n", w.Url, w.LatencyData.Avg, 0, 0, w.SizeData.Avg, 0, 0, w.RequestInfo.SuccesCount, w.RequestInfo.RequestCount())
		return
	}

	fmt.Printf("%s: Latency - Avg: %.2f, Max: %d, Min: %d Size - Avg: %.2f, Max: %d, Min: %d %d/%d \n", w.Url, w.LatencyData.Avg, w.LatencyData.Max, w.LatencyData.Min, w.SizeData.Avg, w.SizeData.Max, w.SizeData.Min, w.RequestInfo.SuccesCount, w.RequestInfo.RequestCount())
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

func (w *Website) MeasureRequestAsync(done chan<-bool) {
	w.MeasureRequest()
	done <- true
}

func (w *Website) MeasureRequest() *Website {
	w.mu.Lock()
	defer w.mu.Unlock()
	defer w.Print()

	startTime := time.Now()
	resp, err := w.httpClient.Get(w.Url)

	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 400 {
		w.RequestInfo.FailCount++
		return w
	}

	measuredTime := int(time.Since(startTime).Milliseconds())
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		w.RequestInfo.FailCount++
		return w
	}

	w.addMeasurements(len(body), measuredTime)
	return w
}
