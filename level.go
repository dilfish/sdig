// sean at shanghai

package main

import "log"

const RootName = "all"

// Level defines a district level
// when a level has members
// they do not need to have a useip
// on the other side, when a level does
// not have a level, ie, it's the lowest
// level, it should have a useip
type Level struct {
	Name string
	UseIP string
	Level int
	Ret RetInfo
	Member []*Level
}


func PrintLevel(level *Level) {
	log.Println("level", level.Level, level.Name, " has ", len(level.Member), "memebers")
	for _, v := range level.Member {
		log.Println("member name is:", v.Name)
	}
	for _, v := range level.Member {
		PrintLevel(v)
	}
}

func CompressLevel(m map[string]string, level *Level, ipMap map[string]string) *Level {
	var children []string
	// pick up all direct children
	for k, v := range m {
		if v == level.Name {
			children = append(children, k)
		}
	}
	// delete children from map
	for _, c := range children {
		delete(m, c)
	}
	// fill children
	for _, c := range children {
		var cl Level
		cl.Name = c
		cl.Level = level.Level + 1
		CompressLevel(m, &cl, ipMap)
		level.Member = append(level.Member, &cl)
	}
	level.UseIP = ipMap[level.Name]
	return level
}

// InitLevel fills the map with relative
// of regions
func InitLevel(ipMap map[string]string) *Level {
	mpd := make(map[string]string)
	// province level
	mpd["辽宁"] = "东北"
	mpd["吉林"] = "东北"
	mpd["黑龙江"] = "东北"
	mpd["河北"] = "华北"
	mpd["山西"] = "华北"
	mpd["内蒙古"] = "华北"
	mpd["北京"] = "华北"
	mpd["天津"] = "华北"
	mpd["山东"] = "华东"
	mpd["江苏"] = "华东"
	mpd["安徽"] = "华东"
	mpd["浙江"] = "华东"
	mpd["福建"] = "华东"
	mpd["江西"] = "华东"
	mpd["上海"] = "华东"
	mpd["湖南"] = "华中"
	mpd["湖北"] = "华中"
	mpd["河南"] = "华中"
	mpd["广东"] = "华南"
	mpd["广西"] = "华南"
	mpd["海南"] = "华南"
	mpd["云南"] = "西南"
	mpd["贵州"] = "西南"
	mpd["四川"] = "西南"
	mpd["西藏"] = "西南"
	mpd["重庆"] = "西南"
	mpd["新疆"] = "西北"
	mpd["陕西"] = "西北"
	mpd["宁夏"] = "西北"
	mpd["青海"] = "西北"
	mpd["甘肃"] = "西北"
	// district level
	mpd["东北"] = "全国"
	mpd["华北"] = "全国"
	mpd["华东"] = "全国"
	mpd["华中"] = "全国"
	mpd["华南"] = "全国"
	mpd["西南"] = "全国"
	mpd["西北"] = "全国"
	// isp level
	mpi := make(map[string]string)
	mpi["电信"] = "1"
	mpi["移动"] = "1"
	mpi["联通"] = "1"
	mpi["教育网"] = "1"
	// multiplex
	mpm := make(map[string]string)
	for k, v := range mpd {
		for kk, _ := range mpi {
			mpm[k + kk] = v + kk
		}
	}
	for k, _ := range mpi {
		mpm["全国" + k] = "大中华区"
	}
	mpm["香港"] = "港澳台"
	mpm["澳门"] = "港澳台"
	mpm["台湾"] = "港澳台"
	mpm["港澳台"] = "大中华区"
	mpm["大中华区"] = "亚洲"
	mpm["日本"] = "亚洲"
	mpm["亚洲"] = RootName
	var level Level
	level.Name = RootName
	level.Level = 1
	// no ip for RootName
	return CompressLevel(mpm, &level, ipMap)
}
