// sean at shanghai

package main

import (
	"errors"
	"log"
	"net"
	"strings"

	"github.com/dilfish/tools"
)

var ipMap map[string]string

func cb(line string) error {
	// comment
	log.Println("read mmap", line)
	if len(line) > 0 && line[0] == '#' {
		return nil
	}
	// space
	if len(line) < 2 {
		return nil
	}
	arr := strings.Split(line, " ")
	if len(arr) != 2 {
		log.Println("format is not 2", arr)
		return errors.New("bad map format")
	}
	if net.ParseIP(arr[1]) == nil {
		log.Println("format is not ip:" + arr[1])
		return errors.New("bad ip format")
	}
	ipMap[arr[0]] = arr[1]
	log.Println("read map", arr[0], arr[1])
	return nil
}

// NewMap read geo-ip information from a txt file
func NewMap(path string) map[string]string {
	ipMap = make(map[string]string)
	err := tools.ReadLine(path, cb)
	if err != nil {
		log.Println("read line error:", err)
		return nil
	}
	return ipMap
}
