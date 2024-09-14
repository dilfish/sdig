// sean at shanghai

package main

import (
	"flag"
	"log"
	"net"

	"github.com/miekg/dns"
)

const DefaultMapFn = ".tmp.sdig.txt"

var flagDomain = flag.String("d", "dilfish.dev.", "domain name")

// 8.8.8.8 support edns
var flagServer = flag.String("s", "119.29.29.29", "server ip")
var flagMapPath = flag.String("m", "/tmp/map.txt", "map file path")
var flagType = flag.String("t", "A", "dns type, default A")
var flagVerbose = flag.Bool("v", false, "verbose print")
var flagRateLimit = flag.Int("r", 1000, "rate limit qps")

// rate limiter
var rateLimiter *QPSRateLimiter

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
	if domain[len(domain)-1] != '.' {
		domain = domain + "."
	}
	if *flagRateLimit != 0 {
		rateLimiter = NewQPSRateLimiter(*flagRateLimit)
	}
	mp := NewMap(*flagMapPath)
	if len(mp) == 0 {
		log.Println("map file does not exist, try to use default one: " + DefaultMapFn)
		mp = NewMap(DefaultMapFn)
	}
	if len(mp) == 0 {
		log.Println("no map file could be used, using built-in map file: " + DefaultMapFn)
		err := GenDefaultMap(DefaultMapFn)
		if err != nil {
			log.Println("gen default map error:", err)
			return
		}
		mp = NewMap(DefaultMapFn)
		if len(mp) == 0 {
			log.Println("default map error")
			return
		}
	}
	var req RequestArgs
	req.Type = uint16(t)
	req.Domain = domain
	for k, v := range mp {
		req.View = k
		req.ViewIP = v
		RunRequest(req)
	}
	PrintResult()
}
