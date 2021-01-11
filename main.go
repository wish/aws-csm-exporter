package main

import (
	"encoding/json"
	"flag"
	"net"

	"github.com/sirupsen/logrus"
)

var (
	csmPort   = flag.Int("csm-port", 1234, "AWS CSM port number")
	csmHost   = flag.String("csm-host", "localhost", "AWS CSM host")
	servePort = flag.Int("serve-port", 8080, "Port to serve metrics")
)

func listenForPackets() {
	addr := net.UDPAddr{
		Port: *csmPort,
		IP:   net.ParseIP(*csmHost),
	}
	connection, err := net.ListenUDP("udp", &addr)
	if err != nil {
		logrus.Infof("Error listening to CSM address: %v", err)
	}
	logrus.Infof("Listening on %s:%d", *csmHost, *csmPort)
	var buf [1024]byte
	for {
		size, _, err := connection.ReadFromUDP(buf[:])
		if err != nil {
			logrus.Errorf("Error receiving packet: %v", err)
		} else {
			logrus.Infof("Received Packet: %s", buf[0:size])
			data := &AWSMetricsData{}
			json.Unmarshal(buf[0:size], data)
			recordMetric(data)
		}
	}
}

func main() {
	flag.Parse()
	registerPrometheusMetrics()
	go listenForPackets()
	serveMetrics(servePort)
}
