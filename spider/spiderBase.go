package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	filepath2 "path/filepath"
	"strings"
	"time"
)
var headers = map[string]string{
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36",
}

func getBaidu() {
	url := "https://www.baidu.com/"
	client := &http.Client{}
	req,_ := http.NewRequest("Get", url, nil)
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(body))
}

func getBiliImgs(url string) []string  {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// 下面使用 goquery 解析网页
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// 用Find()
	// *****************练 习*************************
	// 这部分是练习，熟悉一下goquery
	// 先试试解析头部信息
	// selector 跟其他选择器类似
	// "tag" ".class" "#id"
	head := doc.Find("head")
	// 获取编码
	charsetMeta := head.Find("meta[charset]")
	// Attr()用来获取参数，返回两个值，如果没找到就返回false
	charset,_ := charsetMeta.Attr("charset")
	fmt.Println(charset)
	// 获取title
	title := doc.Find("title").Text()
	fmt.Println(title)

	// 获取body里面的title
	bodyTitle := doc.Find(".title-container h1").Text()
	fmt.Println(bodyTitle)
	// 以上是针对一个元素的情况，下面就来获取图片

	// ******************获 取 图 片************
	var imgs []string
	doc.Find(".article-holder figure").Each(func(i int, s *goquery.Selection){
		img, exits := s.Find("img").Attr("data-src")
		if exits {
			//fmt.Println(img)
			imgs = append(imgs, "https:"+img)
		}
	})
	fmt.Println("图片数量为：", len(imgs))
	fmt.Println("第一张图片的地址为：\n",imgs[0])
	//for _,v := range imgs {
	//	fmt.Println(v)
	//}
	return imgs
}

func fileExists(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// 新建一个带缓冲的channel
var c = make(chan int, 100)

// 一次下载一张图片
func downloadImgs(url string, folder string){
	names := strings.Split(url, "/")
	// 获取文件名
	name := names[len(names)-1]
	path := filepath2.Join(folder, name)
	if fileExists(path) == false {
		fmt.Println(path)
		out, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		pix, err := ioutil.ReadAll(res.Body)
		fmt.Println("downloading: ", url)
		_, err = io.Copy(out, bytes.NewReader(pix))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println(path, "exists")
	}
	c <- 1
}



func main() {
	//getBaidu()
	url := "https://www.bilibili.com/read/cv7586334?from=search"
	imgs := getBiliImgs(url)
	//downloadImgs(imgs[2])

	folder := "bili"
	if fileExists(folder) == false {
		_ = os.Mkdir(folder, os.ModePerm)
	}

	// 顺序下载并计算时间
	//t1 := time.Now()  //获取本地现在时间
	//for _,v := range imgs {
	//	downloadImgs(v, folder)
	//}
	//t2 := time.Now()
	//d := t2.Sub(t1)  //两个时间相减
	//fmt.Println("下载耗时：",d)

	// 试试并行下载
	// 现在对Go的并行还不太熟， c <- 1 到底应该放在函数的哪里
	t1 := time.Now()  //获取本地现在时间
	for _, v := range imgs {
		go downloadImgs(v, folder)
	}
	for i:=0; i< len(imgs); i++ {
		<- c
	}
	t2 := time.Now()
	d := t2.Sub(t1)  //两个时间相减
	fmt.Println("下载耗时：",d)
}
