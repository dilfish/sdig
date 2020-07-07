package main


import (
	"errors"
	"fmt"
	"github.com/dilfish/tools"
	"io/ioutil"
	"net"
	"net/http"
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
	err := tools.ReadLine(*flagMapPath, cb)
	if err != nil {
		fmt.Println("read line error:", err)
		return nil
	}
	return ipMap
}

func GetIPInfo(checkUrl, ip string) string {
	resp, err := http.Get(checkUrl + "/" + ip)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}
	arr := strings.Split(string(bt), "-")
	if len(arr) != 6 {
		return "bad info:" + string(bt)
	}
	return arr[1] + arr[3]
}

func SelfCheck(checkUrl string) error {
	mp := NewMap()
	if len(mp) == 0 {
		return errors.New("bad map init")
	}
	for k, v := range mp {
		ret := GetIPInfo(checkUrl, v)
		if k != ret {
			return errors.New("bad ip info:" + v +":" +  k+ ", it should be:" + ret)
		}
	}
	return nil
}
