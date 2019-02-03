package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const host = "http://www.resgain.net"

var xinSerRe = regexp.MustCompile(`<li ><a href="(/xsdq_[\w].html)">[\w]</a></li>`)
var xinRe = regexp.MustCompile(`<a class="btn btn2" href="//([\w]+.resgain.net)"[^>]+>[^<]+</a>`)

var xinDetail = regexp.MustCompile(`<title>([^<]+)</title>`)
var nameRe = regexp.MustCompile(`<a href="/name/([^.]+).html" class="btn btn-link" target="_blank">[^<]+</a>`)

var done = make(chan int)
var collection *mongo.Collection

func initMongo() {
	client, err := mongo.Connect(context.TODO(), "mongodb://localhost:27017")

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database("girlName").Collection("song")
}

func main() {
	initMongo()
	xinshi := "http://www.resgain.net/xsdq.html"
	fetch(xinshi, func(bytes []byte) {
		matches := xinSerRe.FindAllSubmatch(bytes, -1)
		for _, m := range matches {
			go handleXinSer(host + string(m[1]))
		}
	})
	<-done
}

func handleXinSer(url string) {
	fetch(url, func(body []byte) {
		matches := xinRe.FindAllSubmatch(body, -1)
		for _, m := range matches {
			handleXinList("http://" + string(m[1]) + "/name/girls_%d.html")
		}
	})
}

func handleXinList(url string) {
	for i := 1; i <= 10; i++ {
		handleItemXinGirl(fmt.Sprintf(url, i))
	}
}

func handleItemXinGirl(url string) {
	fetch(url, func(body []byte) {
		matche := xinDetail.FindSubmatch(body)
		xinIndex := strings.Index(string(matche[1]), "姓")
		var xin string
		if xinIndex >= 0 {
			xin = string(matche[1][:3])
		} else {
			xin = string(matche[1][:6])
		}
		matches := nameRe.FindAllSubmatch(body, -1)
		for _, m := range matches {
			name := "宋" + strings.TrimLeft(string(m[1]), xin)
			handleName(name)
		}
	})
}

var nameDetail = regexp.MustCompile(`girls_[\d]+`)

type Girl struct {
	Name       string
	Count      int
	BoyMatch   string
	GirlMatch  string
	Explain    string
	Verses     string
	Five       string
	Three      string
	Tian       int
	Di         int
	Ren        int
	Zong       int
	Wai        int
	TianDetail string
	DiDetail   string
	RenDetail  string
	ZongDetail string
	WaiDetail  string
}

var countReg = regexp.MustCompile(`已有([\d]+)人叫[^<]`)
var boyReg = regexp.MustCompile(`([^>]+)%情况下适用于男孩名`)
var girlReg = regexp.MustCompile(`([^>]+)%情况下适用于女孩名`)
var explainReg = regexp.MustCompile(`<div class="panel-body">[^<]*<strong>([^<]*)</strong>[^<]*</div>`)
var versesReg = regexp.MustCompile(`<h4>([^<]+)</h4>`)
var fiveReg = regexp.MustCompile(`<strong class="titlt">名字五行：</strong>[^<]+<blockquote>([^<]+)</blockquote>`)
var threeReg = regexp.MustCompile(`<strong class="titlt">三才配置：</strong>[^<]+<blockquote>([^<]+)</blockquote>`)

var tianReg = regexp.MustCompile(`<strong>天格</strong>：([\d]+)`)
var diReg = regexp.MustCompile(`<strong>地格</strong>：([\d]+)`)
var renReg = regexp.MustCompile(`<strong>人格</strong>：([\d]+)`)
var zongReg = regexp.MustCompile(`<strong>总格</strong>：([\d]+)`)
var waiReg = regexp.MustCompile(`<strong>外格</strong>：([\d]+)`)

var tianDetailReg = regexp.MustCompile(`<div><strong>天格</strong>：([^<]*)</div>`)
var renDetailReg = regexp.MustCompile(`<div><strong>人格</strong>：([^<]*)</div>`)
var diDetailReg = regexp.MustCompile(`<div><strong>地格</strong>：([^<]*)</div>`)
var waiDetailReg = regexp.MustCompile(`<div><strong>外格</strong>：([^<]*)</div>`)
var zongDetailReg = regexp.MustCompile(`<div><strong>总格</strong>：([^<]*)</div>`)

var unique = make(map[string]byte)

var lock = sync.RWMutex{}

func handleName(name string) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := unique[name]; ok {
		return
	} else {
		unique[name] = 1
	}

	url := "http://song.resgain.net/name/" + name + ".html"
	fetch(url, func(body []byte) {
		g := &Girl{}
		g.Name = name
		match := countReg.FindSubmatch(body)
		if len(match) > 1 {
			count, _ := strconv.Atoi(string(match[1]))
			g.Count = count
		}

		match = boyReg.FindSubmatch(body)
		if len(match) > 1 {
			g.BoyMatch = string(match[1])
		}

		match = girlReg.FindSubmatch(body)
		if len(match) > 1 {
			g.GirlMatch = string(match[1])
		}

		match = explainReg.FindSubmatch(body)
		if len(match) > 1 {
			g.Explain = string(match[1])
		}

		matches := versesReg.FindAllSubmatch(body, -1)
		if len(matches) > 1 {
			length := len(matches)
			verser := ""
			for i := 0; i < length-1; i++ {
				m := matches[i]
				verser += string(m[1])
				verser += "\n"
			}
			g.Verses = strings.TrimRight(verser, "\n")
		}

		match = fiveReg.FindSubmatch(body)
		if len(match) > 1 {
			g.Five = string(match[1])
		}

		match = threeReg.FindSubmatch(body)
		if len(match) > 1 {
			g.Three = string(match[1])
		}

		match = tianReg.FindSubmatch(body)
		if len(match) > 1 {
			num, _ := strconv.Atoi(string(match[1]))
			g.Tian = num
		}

		match = diReg.FindSubmatch(body)
		if len(match) > 1 {

			num, _ := strconv.Atoi(string(match[1]))
			g.Tian = num
		}

		match = renReg.FindSubmatch(body)
		if len(match) > 1 {

			num, _ := strconv.Atoi(string(match[1]))
			g.Tian = num
		}

		match = zongReg.FindSubmatch(body)
		if len(match) > 1 {
			num, _ := strconv.Atoi(string(match[1]))
			g.Tian = num
		}

		match = waiReg.FindSubmatch(body)
		if len(match) > 1 {

			num, _ := strconv.Atoi(string(match[1]))
			g.Tian = num
		}

		match = tianDetailReg.FindSubmatch(body)
		if len(match) > 1 {
			g.TianDetail = string(match[1])
		}

		match = renDetailReg.FindSubmatch(body)
		if len(match) > 1 {
			g.RenDetail = string(match[1])
		}

		match = diDetailReg.FindSubmatch(body)
		if len(match) > 1 {
			g.DiDetail = string(match[1])
		}

		match = waiDetailReg.FindSubmatch(body)
		if len(match) > 1 {
			g.WaiDetail = string(match[1])
		}

		match = zongDetailReg.FindSubmatch(body)
		if len(match) > 1 {
			g.ZongDetail = string(match[1])
		}
		insertResult, err := collection.InsertOne(context.TODO(), g)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	})
}

func fetch(url string, cb func([]byte)) {
	fmt.Println("fetch url", url)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	cb(body)
}
