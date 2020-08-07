package main

import (
	"flag"
	"log"
	"net"
)

var flagDomain = flag.String("d", "www.qiniu.com", "domain name")
var flagServer = flag.String("s", "119.29.29.29", "server ip")
var flagMapPath = flag.String("m", "/tmp/map.txt", "map file path")


func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if net.ParseIP(*flagServer) == nil {
		log.Println("bad server ip", *flagServer)
		return
	}
	if *flagDomain == "" {
		log.Println("empty domain")
		return
	}
	mp := NewMap(*flagMapPath)
	if len(mp) == 0 {
		log.Println("map file does not exist")
		return
	}
	level := InitLevel(mp)
	finish := make(chan struct{})
	go RunLevel(*flagDomain, level, finish)
	<-finish
	close(finish)
	PrintRet(level, "")
	return
}
