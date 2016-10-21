package main

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/astaxie/beego"
)

//--------------------ddns
type myRegexp struct{ *regexp.Regexp }

//add a new method to our new regular expression type
func (r *myRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)
	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}
	for i, name := range r.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]
	}
	return captures
}

func getIPAddress(str string) (string, bool) {
	myExp := myRegexp{regexp.MustCompile(`var url='http://(?P<first>\d+)\.(?P<second>\d+).(?P<third>\d+).(?P<fourth>\d+):`)}
	mmap := myExp.FindStringSubmatchMap(str)
	ip := mmap["first"] + "." + mmap["second"] + "." + mmap["third"] + "." + mmap["fourth"]
	return ip, len(mmap) == 4
}

func getDdns(uri string) (string, bool) {
	myExp := myRegexp{regexp.MustCompile(`/hls/(?P<first>[\w]+|[\d]+)/`)}
	mmap := myExp.FindStringSubmatchMap(uri)
	return mmap["first"], len(mmap) == 1
}

func getDdnsIP(dns string) (string, bool) {
	resp, err := http.Get("http://www.hik-online.com/" + dns)
	if err != nil {
		return "", false
	}

	defer resp.Body.Close()
	br := bufio.NewReader(resp.Body)
	line, isPrefix, err := br.ReadLine()
	for i := 1; err == nil && !isPrefix; i++ {
		if i == 113 {
			s := string(line)
			ip, err := getIPAddress(s)
			if err == true {
				return ip, true
			}
			return "", false
		}
		line, isPrefix, err = br.ReadLine()
	}

	return "", false
}

//--------------------

type BeeferController struct {
	beego.Controller
}

func (c *BeeferController) Get() {
	// 断断这个流是否已经启动ffmpeg去拉了

	// 启动一个exec ffmpeg 去拉流并推到nginx-rtmp

	uri := c.Ctx.Request.URL.String()
	if dns, err := getDdns(uri); err == true {
		if ip, err := getDdnsIP(dns); err == true {
			fmt.Println(ip)
			m3u8 := "#EXTM3U\r\n#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=200000\r\nhttp://" + ip + uri
			c.Ctx.WriteString(m3u8)
			return
		}
	}

	//TODO Return 404
	fmt.Println(uri)
}

var c1 = make(chan string)
var c2 = make(chan string)

func startFfmpeg1(uri string, c chan string) {
	time.Sleep(time.Second)
	c <- uri
}

func startFfmpeg2(uri string, c chan string) {
	time.Sleep(time.Second * 2)
	c <- uri
}

func checkFfmpeg(m map[string]bool, c chan string) {
	for {
		select {
		case msg := <-c:
			m[msg] = false
			fmt.Println("received", msg)
		}
	}
}

func main() {
	c := make(chan string)
	m := make(map[string]bool)
	uri1 := "hls/tnwl123456/ch34.m3u8"

	if _, ok := m[uri1]; !ok { // 如果不存在这个key，就产生它
		m[uri1] = false
	}

	go checkFfmpeg(m, c)
	if m[uri1] == false { // 如果这个流没有在拉，就启动拉流
		m[uri1] = true
		// TODO
		go startFfmpeg1(uri1, c)
	}
	time.Sleep(time.Second * 2)

	if m[uri1] == false { // 如果这个流没有在拉，就启动拉流
		m[uri1] = true
		// TODO
		go startFfmpeg1(uri1, c)
	}

	beego.Router("/hls/*", &BeeferController{})
	beego.Run(":8000")
}
