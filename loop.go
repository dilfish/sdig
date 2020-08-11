// sean at shanghai

package main

import "log"

// SameResult checkes if this array has same result
// if so, return the one
func SameResult(member []*Level) string {
    var same string
    for _, m := range member {
        if same == "" {
            same = m.Ret.Result
        }
        // not same
        if m.Ret.Result != same {
            return ""
        }
    }
    return same
}

// RunLevel run dns request level by level
func RunLevel(req RequestArgs, level *Level, finish chan struct{}) {
	children := make(chan struct{})
	for _, m := range level.Member {
		go RunLevel(req, m, children)
	}
	if level.UseIP != "" {
		ch := make(chan RetInfo)
		go MakeRequest(req.Domain, level.UseIP, level.Name, req.Type, ch)
		level.Ret = <-ch
		close(ch)
	}
	for _, _ = range level.Member {
		<-children
	}
	close(children)
        same := SameResult(level.Member)
        if same != "" {
            level.Ret.Result = same
        }
	finish<-struct{}{}
	return
}

// PrintRet print all result by level
func PrintRet(level *Level, parentResult string) {
	if level.Ret.Result != parentResult && level.Ret.Result != "" {
		log.Println("result:", level.Name)
		log.Println("\n" + level.Ret.Result)
	}
	for _, m := range level.Member {
		if level.Ret.Result != "" {
			parentResult = level.Ret.Result
		}
		PrintRet(m, parentResult)
	}
}
