// sean at shanghai

package main

import "log"

func RunLevel(domain string, level *Level, finish chan struct{}) {
	children := make(chan struct{})
	for _, m := range level.Member {
		go RunLevel(domain, m, children)
	}
	if level.UseIP != "" {
		ch := make(chan RetInfo)
		go MakeRequest(domain, level.UseIP, level.Name, ch)
		level.Ret = <-ch
		close(ch)
	}
	for _, _ = range level.Member {
		<-children
	}
	close(children)
	finish<-struct{}{}
	return
}

func PrintRet(level *Level, parentResult string) {
	if level.Ret.Result != parentResult && level.Ret.Result != "" {
		log.Println("result", level.Name, level.Ret.Result)
	}
	for _, m := range level.Member {
		if level.Ret.Result != "" {
			parentResult = level.Ret.Result
		}
		PrintRet(m, parentResult)
	}
}