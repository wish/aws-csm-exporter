package main

type AWSMetricsData struct {
	ClientID                 string
	API                      string
	Service                  string
	Timestamp                int64
	Type                     string
	Version                  int
	AttemptCount             int
	Latency                  int
	Fqdn                     string
	UserAgent                string
	AttemptLatency           int
	SessionToken             string
	Region                   string
	AccessKey                string
	HTTPStatusCode           int
	XAmzID2                  string
	XAmzRequestID            string
	AWSException             string
	AWSExceptionMessage      string
	SDKException             string
	FinalHTTPStatusCode      int
	FinalAWSException        string
	FinalAWSExceptionMessage string
	FinalSDKException        string
	FinalSDKExceptionMessage string
	DestinationIP            string
	ConnectionReused         int
	AcquireConnectionLatency int
	ConnectLatency           int
	RequestLatency           int
	DNSLatency               int
	TCPLatency               int
	SSLLatency               int
	MaxRetriesExceeded       int
}
