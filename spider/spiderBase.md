# Go简单爬虫

## 1.0 获取baidu首页

```go
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
    // 奇怪的是，如果只用 http.get()，响应就是错误的页面
    // 如果按上面来获取bibili首页，则是500错误
    // 直接用http.get()就行
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
```



## 1.1 使用goquery解析网页

导入goquery

```go
go get github.com/PuerkitoBio/goquery
```



就拿获取B站专栏图片为例



```go
func spider() {
	res, err := http.Get("https://www.bilibili.com/read/cv7533939")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	head := doc.Find("head")

	// 添加一个字符串数组（或者一个map）来存imgs
	var imgs = []string{}

	// 试试打印图片地址
	doc.Find(".article-holder figure").Each(func(i int, s *goquery.Selection){
		img, exits := s.Find("img").Attr("data-src")
		if exits {
			//fmt.Println(img)
			imgs = append(imgs, "https:"+img)
		}
	})
	fmt.Println(imgs[1])
	for _,v := range imgs {
		fmt.Println(v)
	}
	fmt.Println(len(imgs))
}
```

