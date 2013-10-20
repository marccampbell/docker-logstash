package main

import (
    "net"
    "fmt"
    "bufio"
    "encoding/json"
    "time"
)

type dockerWatcher struct {
    dockerConn net.Conn
    logstashAddr *net.TCPAddr
}

type DockerEvent struct {
    Status string
    Id string
    From string
    Time int64
}

type LogstashEvent struct {
    Timestamp string `json:"@timestamp"`
    Tags []string `json:"tags"`
    Message string `json:"message"`
    ContainerID string `json:"container_id"`
    Event string `json:"event"`
}

var e *dockerWatcher

func newDockerWatcher() *dockerWatcher {
    dockerConn, _ := net.Dial("unix", "/var/run/docker.sock")
    logstashAddr, _ := net.ResolveTCPAddr("tcp", "172.17.42.1:9125")

    return &dockerWatcher{
        dockerConn: dockerConn,
        logstashAddr: logstashAddr,
    }
}

func (w *dockerWatcher) Listen() {
    fmt.Fprintf(w.dockerConn, "GET /v1.6/events HTTP/1.1\r\n\r\n")
    r := bufio.NewReader(w.dockerConn)

    for {
        body, err := r.ReadString('\n')
        if err != nil {
            fmt.Println("error in readstring from events feed")
            return
        } else {
            logstashConn, _ := net.DialTCP("tcp", nil, w.logstashAddr)
            defer logstashConn.Close()
            sendContainerEvent := func(dockerEvent *DockerEvent) {
                logstashEvent := &LogstashEvent{Timestamp: time.Unix(dockerEvent.Time, 0).Format("2006-01-02T15:04:05Z07:00"),
                                                Tags: []string{"docker", "docker_" + dockerEvent.Status, dockerEvent.From},
                                                ContainerID: dockerEvent.Id,
                                                Event: dockerEvent.Status}
                if dockerEvent.Status == "create" {
                    logstashEvent.Message = "created a docker container from image " + dockerEvent.From
                } else if dockerEvent.Status == "start" {
                    logstashEvent.Message = "started a docker container"
                } else if dockerEvent.Status == "die" {
                    logstashEvent.Message = "docker container died"
                } else {
                    logstashEvent.Message = "other docker event"
                }

                b, err := json.Marshal(logstashEvent)
                if err != nil {
                    fmt.Println("json error:", err)
                } else {
                    logstashConn.Write([]byte(b))
                }
            }

            var dockerEvent DockerEvent
            err := json.Unmarshal([]byte(body), &dockerEvent)
            if err == nil {
                go sendContainerEvent(&dockerEvent)
            }
        }
    }
}

