package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"sort"
	"strings"
	"time"
)

var flagDomain = flag.String("d", "www.qiniu.com", "domain name")
var flagServer = flag.String("s", "119.29.29.29", "server ip")

func NewRequest(domain, clientIP string) *dns.Msg {
	req := new(dns.Msg)
	req.SetQuestion(domain, dns.TypeA)
	if clientIP != "" {
		ip := net.ParseIP(clientIP)
		if ip == nil {
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

func MakeCall(req *dns.Msg) (msg *dns.Msg, duration time.Duration, err error) {
	c := new(dns.Client)
	s := *flagServer + ":53"
	return c.Exchange(req, s)
}


type RetInfo struct {
	Region string
	Result string
}


func MakeRequest(domain, clientIP, region string, ch chan RetInfo) {
	var r RetInfo
	r.Region = region
	req := NewRequest(domain, clientIP)
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


const TypeMob = 1
const TypeTel = 2
const TypeCnc = 3
const TypeEdu = 4
const TypeOther = 5
func GetType(str string) int {
	if strings.Contains(str, "移动") {
		return TypeMob
	}
	if strings.Contains(str, "电信") {
		return TypeTel
	}
	if strings.Contains(str, "联通") {
		return TypeCnc
	}
	if strings.Contains(str, "教育网") {
		return TypeEdu
	}
	return TypeOther
}


func CallWithIPList(domain string) error {
	if domain[len(domain) - 1] != '.' {
		domain = domain + "."
	}
	mp := NewMap()
	if mp == nil {
		fmt.Println("new map error")
		return errors.New("read ip file error")
	}
	ch := make(chan RetInfo)
	var  n int
	for k, v := range mp {
		go MakeRequest(domain, v, k, ch)
		n = n + 1
	}
	mpTel := make(map[string][]string)
	mpCnc := make(map[string][]string)
	mpMob := make(map[string][]string)
	mpEdu := make(map[string][]string)
	mpOther := make(map[string][]string)
	for i := 0;i < n;i ++ {
		r := <-ch
		tp := GetType(r.Region)
		switch tp {
		case TypeTel:
			v := mpTel[r.Result]
			v = append(v, r.Region)
			mpTel[r.Result] = v
		case TypeMob:
			v := mpMob[r.Result]
			v = append(v, r.Region)
			mpMob[r.Result] = v
		case TypeCnc:
			v := mpCnc[r.Result]
			v = append(v, r.Region)
			mpCnc[r.Result] = v
		case TypeEdu:
			v := mpEdu[r.Result]
			v = append(v, r.Region)
			mpEdu[r.Result] = v
		default:
			v := mpOther[r.Result]
			v = append(v, r.Region)
			mpOther[r.Result] = v
		}
	}
	close(ch)
	fmt.Println("----------------------")
	for k, v := range mpCnc {
		for _, ip := range v {
			fmt.Println(ip)
		}
		fmt.Println(k)
	}
	fmt.Println("----------------------")
	for k, v := range mpMob {
		for _, ip := range v {
			fmt.Println(ip)
		}
		fmt.Println(k)
	}
	fmt.Println("----------------------")
	for k, v := range mpTel {
		for _, ip := range v {
			fmt.Println(ip)
		}
		fmt.Println(k)
	}
	fmt.Println("----------------------")
	for k, v := range mpEdu {
		for _, ip := range v {
			fmt.Println(ip)
		}
		fmt.Println(k)
	}
	fmt.Println("----------------------")
	for k, v := range mpOther {
		for _, ip := range v {
			fmt.Println(ip)
		}
		fmt.Println(k)
	}
	return nil
}

func main() {
	flag.Parse()
	if net.ParseIP(*flagServer) == nil {
		log.Println("bad server ip", *flagServer)
		return
	}
	err := CallWithIPList(*flagDomain)
	if err != nil {
		fmt.Println("call with ip list error:", err)
	}
	return
}
