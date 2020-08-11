// sean at shanghai

package main

import (
	"github.com/miekg/dns"
	"log"
	"net"
	"sort"
	"strings"
	"time"
)

type RequestArgs struct {
	Domain string
	Type uint16
}

// NewRequest make a dns request struct with specified domain
// and client ip
func NewRequest(domain, clientIP string, tp uint16) *dns.Msg {
	req := new(dns.Msg)
	req.SetQuestion(domain, tp)
	if clientIP != "" {
		ip := net.ParseIP(clientIP)
		if ip == nil {
			log.Println("bad ip", clientIP)
			return req
		}
		opt := new(dns.OPT)
		opt.Hdr.Name = "."
		opt.Hdr.Rrtype = dns.TypeOPT
		e := new(dns.EDNS0_SUBNET)
		e.Code = dns.EDNS0SUBNET
		e.Family = 1
		e.SourceNetmask = 32
		e.SourceScope = 0
		e.Address = ip
		opt.Option = append(opt.Option, e)
		req.Extra = []dns.RR{opt}
	}
	return req
}

// MakeCall calls the request
func MakeCall(req *dns.Msg) (msg *dns.Msg, duration time.Duration, err error) {
	c := new(dns.Client)
	s := *flagServer + ":53"
	count := 0
	for {
		msg, duration, err = c.Exchange(req, s)
		if err != nil && strings.Index(err.Error(), "i/o timeout") >= 0 {
			return
		}
		count = count + 1
		if count > 2 {
			break
		}
	}
	return
}


// RetInfo stands for
// a result of dns request
type RetInfo struct {
	Region string
	Result string
}

// MakeRequest make a request on the fly
// write back result to channel
func MakeRequest(domain, clientIP, region string, tp uint16, ch chan RetInfo) {
	var r RetInfo
	r.Region = region
	req := NewRequest(domain, clientIP, tp)
	msg, _, err := MakeCall(req)
	if err != nil {
		r.Result = err.Error()
	} else {
		strs := make([]string, 0)
		for _, a := range msg.Answer {
			switch a.(type) {
			case *dns.A:
				r := a.(*dns.A)
				strs = append(strs, r.A.String())
			case *dns.CNAME:
				c := a.(*dns.CNAME)
				strs = append(strs, c.Target)
			default:
				strs = append(strs, a.String())
			}
		}
		sort.Strings(strs)
		str := strings.Join(strs, "\n")
		r.Result = str
	}
	ch<-r
}

func IsSupportType(str string) uint16 {
	t := dns.TypeNone
	str = strings.ToUpper(str)
	switch str {
	case "A":
		t = dns.TypeA
	case "AAAA":
		t = dns.TypeAAAA
	}
	return t
}
