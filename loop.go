// sean at shanghai

package main

import (
	"fmt"
	"log"
)

// result -> view list
var retMap map[string][]string

// RunRequest
func RunRequest(req RequestArgs) {
	ret, err := MakeRequest(req.Domain, req.ViewIP, req.View, req.Type)
	if err != nil {
		log.Println("make request error:", req, err)
		return
	}
	if retMap == nil {
		retMap = make(map[string][]string)
	}
	v := retMap[ret.Result]
	v = append(v, ret.Region)
	retMap[ret.Result] = v
}

func PrintResult() {
	for k, v := range retMap {
		for _, view := range v {
			fmt.Println(view + ":")
		}
		fmt.Println(k)
	}
}
