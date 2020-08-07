package main

import "log"

func RunLevel(domain string, level *Level, finish chan struct{}) {
	children := make(chan struct{})
	for _, m := range level.Member {
		go RunLevel(domain, m, children)
	}
	if level.UseIP != "" {
		ch := make(chan RetInfo)
		// log.Println("doing", level.Name)
		go MakeRequest(domain, level.UseIP, level.Name, ch)
		level.Ret = <-ch
		log.Println("I got result", level.Ret)
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
	for _, m := range level.Member {
		p := parentResult
		if level.Ret.Result != "" {
			p = parentResult
		}
		PrintRet(m, p)
	}
	if level.Name == RootName {
		return
	}
	if parentResult != "" && parentResult == level.Ret.Result {
		return
	}
	if level.Ret.Result != "" {
		log.Println("result", level.Name, level.Ret.Result, parentResult)
	}
}