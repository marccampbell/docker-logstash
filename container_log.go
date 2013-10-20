package main

import (
    "net"
)

type containerLog struct {
    dockerConn net.Conn
    logstashAddr *net.TCPAddr
}

type LogstashMessage struct {
    Timestamp string `json:"@timestamp"`
    Tags []string `json:"tags"`
    Message string `json:"message"`
    ContainerID string `json:"container_id"`
}

var c *containerLog

func newContainerLog() *containerLog {
    dockerConn, _ := net.Dial("unix", "/var/run/docker.sock")
    logstashAddr, _ := net.ResolveTCPAddr("udp", "172.17.42.1:9125")

    return &containerLog{
        dockerConn: dockerConn,
        logstashAddr: logstashAddr,
    }
}

func (w *containerLog) Attach() {
}
