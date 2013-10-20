package main

import (
    "fmt"
    "flag"
    "os"
)

var (
    printVersion bool

    verbose bool
    veryVerbose bool
)

const releaseVersion = "0.1"

func init() {
    flag.BoolVar(&printVersion, "version", false, "print the version and exit")

    flag.BoolVar(&verbose, "v", false, "verbose logging")
    flag.BoolVar(&veryVerbose, "vv", false, "very verbose logging")

}

func main() {
    flag.Parse()

    if printVersion {
        fmt.Println(releaseVersion)
        os.Exit(0)
    }

    w := newDockerWatcher()
    w.Listen()

}

