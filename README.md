docker-logstash
===============

The logstash agent for a Docker world

Install this agent on all Docker hosts (not inside containers), and instantly get logstash messages when any of the following events occur:

* Container event (create, start, stop, destroy)

Docker-Logstash is written in Go.

## Getting Started

### Getting docker-logstash

Clone [this repo](https://github.com/marccampbell/docker-logstash) and run `go build`.

### Requirements

Logstash 1.2.1 is required.  A new `input1` is required in logstash:

```
  tcp {
    port => 9125
    type => "docker_container_event"
    codec => "json"
  }
```
