package main

import (
	"github.com/peterq/pan-light/demo/host"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	host.Start()
}
