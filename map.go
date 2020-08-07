package main

import (
	"errors"
	"github.com/dilfish/tools"
	"log"
	"net"
	"strings"
)

var ipMap map[string]string

func cb(line string) error {
	arr := strings.Split(line, " ")
	if len(arr) != 2 {
		log.Println("format is not 2")
		return errors.New("bad map format")
	}
	if net.ParseIP(arr[1]) == nil {
		log.Println("format is not ip:" + arr[1])
		return errors.New("bad ip format")
	}
	ipMap[arr[0]] = arr[1]
	return nil
}

func NewMap(path string) map[string]string {
	ipMap = make(map[string]string)
	err := tools.ReadLine(path, cb)
	if err != nil {
		log.Println("read line error:", err)
		return nil
	}
	return ipMap
}
