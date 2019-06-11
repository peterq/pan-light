package main

import (
	"flag"
	"github.com/peterq/pan-light/server/cmd/cv"
	"github.com/peterq/pan-light/server/cmd/nickname"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Usage = func() {
		println("error usage")
	}
	flag.Parse()
	cmd := "avatar"
	cmd = flag.Arg(0)
	switch cmd {
	case "avatar":
		avatar()
	case "cv":
		cv.CV()
	default:
		flag.Usage()
	}
}

func avatar() {
	nickname.FetchAndSaveAvatarFromInternet()
}
