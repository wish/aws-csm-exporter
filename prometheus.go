package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	awsReqCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "awscsm",
		Name:      "req_count",
		Help:      "Total reqests made to AWS",
	}, []string{"result"})

	awsThrottleCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "awscsm",
		Name:      "throttle_count",
		Help:      "Number of requests that failed due to throttling",
	}, []string{"service"})

	awsReqLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "awscsm",
		Name:      "req_latency",
		Help:      "Latency in ms of aws requests",
		Buckets:   prometheus.LinearBuckets(100, 500, 20),
	}, []string{"service"})
)

func registerPrometheusMetrics() {
	prometheus.MustRegister(awsReqCount)
	prometheus.MustRegister(awsThrottleCount)
	prometheus.MustRegister(awsReqLatency)
	logrus.Infof("Finished registering metrics")
}

func serveMetrics(servePort *int) {
	logrus.Infof("Serving metrics on :%d", *servePort)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", *servePort), nil)
	if err != nil {
		logrus.Errorf("Error serving metrics: %v", err)
	}
}

func checkSuccessCode(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func recordMetric(metricData *AWSMetricsData) {
	switch metricData.Type {
	case "ApiCall":
		if checkSuccessCode(metricData.FinalHTTPStatusCode) {
			awsReqCount.WithLabelValues("success").Add(1)
		} else {
			awsReqCount.WithLabelValues("error").Add(1)
		}
		logrus.Infof("%f", float64(metricData.Latency))
		awsReqLatency.WithLabelValues(metricData.Service).Observe(float64(metricData.Latency))
		break
	case "ApiCallAttempt":
		if checkSuccessCode(metricData.HTTPStatusCode) {
			awsReqCount.WithLabelValues("success").Add(1)
		} else {
			awsReqCount.WithLabelValues("error").Add(1)
		}
		logrus.Infof("%f", float64(metricData.AttemptLatency))
		awsReqLatency.WithLabelValues(metricData.Service).Observe(float64(metricData.AttemptLatency))
		break
	}
	if metricData.FinalHTTPStatusCode == 419 {
		awsThrottleCount.WithLabelValues(metricData.Service).Add(1)
	}
}
