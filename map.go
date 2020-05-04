package main


import (
	"errors"
	"fmt"
	"github.com/dilfish/tools"
	"net"
	"strings"
)

var ipMap map[string]string

func cb(line string) error {
	arr := strings.Split(line, " ")
	if len(arr) != 2 {
		fmt.Println("format is not 2")
		return errors.New("bad map format")
	}
	if net.ParseIP(arr[1]) == nil {
		fmt.Println("format is not ip:" + arr[1])
		return errors.New("bad ip format")
	}
	ipMap[arr[0]] = arr[1]
	return nil
}

func NewMap() map[string]string {
	ipMap = make(map[string]string)
	err := tools.ReadLine("/tmp/map.txt", cb)
	if err != nil {
		fmt.Println("read line error:", err)
		return nil
	}
	return ipMap
}
