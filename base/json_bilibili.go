package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func test() {
	// 1.1 imple
	info := Mine{}
	err := json.Unmarshal([]byte(mine), &info)
	if err != nil{
		log.Fatalln(err)
		//print("error:  ", err)
	}
	fmt.Println(info)
	fmt.Println(info.Name)
	fmt.Println(info.Scores)

	// 1.2 map in map
	info2 := InMap{}
	err = json.Unmarshal([]byte(inmap), &info2)
	if err != nil{
		log.Fatalln(err)
		//print("error:  ", err)
	}
	fmt.Println(info2)
	fmt.Println(info2.Name)
	fmt.Println(info2.Scores)
	fmt.Println(info2.Scores.LessonC)

	var event map[string]interface{}

	err = json.Unmarshal([]byte(inmap), &event)
	if err != nil {
		panic(err)
	}
	fmt.Println(event)
}



var mine = `{
    "name":"Joe",
    "age":23,
    "height":175.5,
    "scores":[80,90,100]
}`

type Mine struct {
	Name 	string 		`json:"name"`
	Age 	int 		`json:"age"`
	Height 	float32 	`json:"height"`
	Scores 	[]int		`json:"scores"`
}

var inmap = `{
	"name":"Joe",
	"age":23,
	"height":175.5,
	"Scores":{
		"C":60,
		"Go":70,
		"Python":80
	}
}`

type InMap struct {
	Name 	string 			`json:"name"`
	Age 	int 			`json:"age"`
	Height 	float32 		`json:"height"`
	Scores 	Scores	`json:"scores"`
}

type Scores struct {
	LessonC 		int 	`json:"C"`
	LessonGo 		int 	`json:"Go"`
	LessonPython 	int 	`json:"Python"`
}

type bili struct {
	Data data  `json:"data"`
}

type data struct {
	Artilce []artilces `json:"articles"`
	Count int `json:"count"`
}

type artilces struct {
	Id 		int 		`json:"id"`
	Title 	string 		`json:"title"`
}


func getCount(url string) int {
	res, err := http.Get(url)
	if err != nil {
		print(err)
	}
	//defer res.Body.Close()
	body, _:= ioutil.ReadAll(res.Body)
	bi := bili{}
	err = json.Unmarshal(body, &bi)
	if err != nil{
		log.Fatalln(err)
	}
	fmt.Println(bi.Data.Count)
	return bi.Data.Count
}

func getids (url string) []int {
	var ids []int
	res, err := http.Get(url)
	if err != nil {
		print(err)
	}
	//defer res.Body.Close()
	body, _:= ioutil.ReadAll(res.Body)
	bi := bili{}
	err = json.Unmarshal(body, &bi)
	if err != nil{
		log.Fatalln(err)
	}
	for _,v :=range bi.Data.Artilce {
		fmt.Println(v.Id, v.Title)
		ids = append(ids, v.Id)
	}
	return ids
}

func getAllids(url string, pages int) []artilces {
	var allIds []artilces
	for i:=1; i< pages+1; i++ {
		currentURL := url + "&pn=" + strconv.Itoa(i) + "&ps=30"
		fmt.Println(currentURL)
		res, err := http.Get(currentURL)
		if err != nil {
			print(err)
		}
		//defer res.Body.Close()
		body, _:= ioutil.ReadAll(res.Body)
		bi := bili{}
		err = json.Unmarshal(body, &bi)
		if err != nil{
			log.Fatalln(err)
		}
		for _,v :=range bi.Data.Artilce {
			var temp artilces
			fmt.Println(v.Id, v.Title)
			temp.Id = v.Id
			temp.Title = v.Title
			allIds = append(allIds, temp)
		}
	}
	return allIds
}


func main() {
	count := getCount("https://api.bilibili.com/x/space/article?mid=18218639&pn=1&ps=1")
	pages := count/30 + 1
	fmt.Println(pages)
	ids := getAllids("https://api.bilibili.com/x/space/article?mid=18218639", pages)
	f,_ := os.Create("ids.txt")
	for _,v := range ids {
		_, _ = f.WriteString(strconv.Itoa(v.Id) + ":" + v.Title + "\n")
	}
	_ = f.Close()
}
