// sean at shanghai
// todo: enhence map.txt
// a better printret function

package main

import (
	"flag"
	"github.com/miekg/dns"
	"log"
	"net"
)

var flagDomain = flag.String("d", "www.qiniu.com", "domain name")
var flagServer = flag.String("s", "119.29.29.29", "server ip")
var flagMapPath = flag.String("m", "/tmp/map.txt", "map file path")
var flagType = flag.String("t", "A", "dns type, default A")


func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if net.ParseIP(*flagServer) == nil {
		log.Println("bad server ip", *flagServer)
		return
	}
	domain := *flagDomain
	if domain == "" {
		log.Println("empty domain")
		return
	}
	t := IsSupportType(*flagType)
	if t == dns.TypeNone {
		log.Println("we do not support type:", *flagType)
		return
	}
	if domain[len(domain) - 1] != '.' {
		domain = domain + "."
	}
	mp := NewMap(*flagMapPath)
	if len(mp) == 0 {
		log.Println("map file does not exist")
		return
	}
	level := InitLevel(mp)
	finish := make(chan struct{})
	var req RequestArgs
	req.Type = uint16(t)
	req.Domain = *flagDomain
	go RunLevel(req, level, finish)
	<-finish
	close(finish)
	PrintRet(level, "")
	return
}
