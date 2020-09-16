package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// b站图片的前缀。因为为了节约空间，存到文件/数据库中的只是文件名，
// 加上前缀才能下载
var perfix = "https://i0.hdslb.com/bfs/article/"

// 下面4个结构体是用来解析json的
// 因为url响应就是多重嵌套的

// 表示一个UP的信息
// 命名有点问题，暂时就这样吧
type bili struct {
	Data data  `json:"data"`
}

// 解析上面的Data
type data struct {
	Artilce []artilces `json:"articles"`
	Count int `json:"count"`
}

// 解析上面的Artilces
type artilces struct {
	Id 		int 		`json:"id"`
	Title 	string 		`json:"title"`
	Author author		`json:"author"`
}

// 获取up的昵称
type author struct {
	Name string `json:"name"`
}

// 下面两个结构体表示需要存储的信息
// 跟上面的几个有些不同，偷偷懒，没想着优化结构

// 一个UP的基本信息
type up struct {
	UID int
	Name string
	Num int
	CVids []int  // 所有专栏的id
	CVs []cv
}

// 专栏信息
type cv struct {
	CVid int
	CVtitle string
	Images []string
}

// 获取up专栏数目和up的昵称
// 输入为up的id int
// 输出为 专栏数 int，昵称 string
func getCountName(id int) (int, string) {
	url := "https://api.bilibili.com/x/space/article?mid="+ strconv.Itoa(id) + "&pn=1&ps=1"
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
	name := bi.Data.Artilce[0].Author.Name

	//fmt.Println(bi.Data.Count)
	return bi.Data.Count, name
}

// 获取所有的专栏id
// 输入：用户id int，专栏数量count，int
// 输出：所有专栏id []int
func getAllids(id int, count int) []int {
	var allIds []int
	url := "https://api.bilibili.com/x/space/article?mid=" + strconv.Itoa(id)
	pages := count/30 + 1
	for i := 1; i < pages + 1; i++ {
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
			//fmt.Println(v.Id, v.Title)
			temp.Id = v.Id
			temp.Title = v.Title
			allIds = append(allIds, temp.Id)
		}
	}
	return allIds
}

// 获取一个页面下的title和所有图片的文件名，因为前缀都是一样的，这样做可减少内存空间
// 输入：专栏的id int
// 输出：专栏title string，所有图片文件名 []string
func getCVImgs(id int) (string, []string) {
	url := "https://www.bilibili.com/read/cv" + strconv.Itoa(id)
	res, err := http.Get(url)
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
	// 添加一个字符串数组（或者一个map）来存imgs
	var imgs []string
	title := doc.Find("title").Text()
	// 试试打印图片地址
	doc.Find(".article-holder figure").Each(func(i int, s *goquery.Selection){
		img, exits := s.Find("img").Attr("data-src")
		if exits {
			imgnames := strings.Split(img, "/")
			name := imgnames[len(imgnames)-1]
			imgs = append(imgs, name)
		}
	})
	//for _,v := range imgs {
	//	fmt.Println(v)
	//}
	return title, imgs
}

// 存为json文件
// 输入：文件名 string， 不带后缀
func saveAsJson(data up, filename string) {
	file, err := json.MarshalIndent(data,"","  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename + ".json", file, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// 存到mongodb
// 输入：up的id
// 这里偷偷懒，连接设置暂时这样，虽然读配置也不难
// 连接设置主要包括：
// 用户名 密码 主机host 端口 数据库 集合
func save(id int) {
	// 首先看看之前有没有存进去
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://localhost")
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// 指定获取要操作的数据集
	filter := bson.D{{"uid", id}}
	var result up
	collection := client.Database("test").Collection("test")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == nil {
		log.Fatal("此用户已保存")
		//os.Exit(1)
	}

	// 如果原来没有就存进去
	uper := up{UID: id}
	uper.Num, uper.Name = getCountName(uper.UID)
	if uper.Num == 0 {
		log.Fatal("没有此用户或该用户无专栏图片")
	}
	uper.CVids = getAllids(uper.UID, uper.Num)
	// 下面就要来循环获取图片了
	for _, v := range uper.CVids {
		fmt.Println("Getting: ", v)
		// 悠一点
		//time.Sleep(time.Second)
		tmp := cv{}
		tmp.CVid = v
		tmp.CVtitle, tmp.Images = getCVImgs(v)
		uper.CVs = append(uper.CVs, tmp)
	}
	//fmt.Println(uper)

	// 存进去
	insertResult, err := collection.InsertOne(context.TODO(), uper)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	// 断开连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
	// 存到文件
	saveAsJson(uper, strconv.Itoa(uper.UID))
}


// 查询
func find(id int) {
	// 下面来读取
	clientOptions := options.Client().ApplyURI("mongodb://localhost")
	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// 指定获取要操作的数据集
	collection := client.Database("test").Collection("test")
	filter := bson.D{{"uid", id}}
	var result up
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	for j, cv :=range result.CVs {
		if j == 1 {
			break
		}
		fmt.Println(cv.CVid)
		for i, imgs := range cv.Images {
			if i == 1 {
				break
			}
			fmt.Println(perfix + imgs)
		}
	}
	saveAsJson(result, strconv.Itoa(result.UID))

	// 断开连接
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func main() {
	// 健壮性不够好
	// 用户id输错
	// 数据库连不上
	// 看看怎么整一个连接池
	var uid = flag.Int("uid", 0, "up主id")
	var search = flag.Bool("find", false, "是否为查询")
	flag.Parse()
	fmt.Println("-uid:", *uid)
	fmt.Println("-search:", *search)
	id := *uid
	//fmt.Println(getCountName(id))
	if id == 0 {
		log.Fatal("请输入正确id")
	}
	if *search {
		find(id)
	} else {
		save(id)
	}
}
