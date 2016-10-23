package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	//"time"

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
	myExp := myRegexp{regexp.MustCompile(`/hls/(?P<first>[\w]+|[\d]+)_`)}
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

var my_c = make(chan string)
var my_m = make(map[string]bool)

func (c *BeeferController) Get() {
	// 启动一个exec ffmpeg 去拉流并推到nginx-rtmp

	uri := c.Ctx.Request.URL.String()
	if dns, err := getDdns(uri); err == true {
		if ip, err := getDdnsIP(dns); err == true {
			// 断断这个流是否已经启动ffmpeg去拉了
			if _, ok := my_m[uri]; !ok { // 如果不存在这个key，就产生它
				my_m[uri] = false
			}

			if my_m[uri] == false { // 如果这个流没有在拉，就启动拉流
				my_m[uri] = true

				input := "rtsp://admin:hk8898878@" + ip + ":554/mpeg4/ch34/sub/av_stream"
				output := "rtmp://42.51.201.196" + uri
				fmt.Println(input, output)
				go startFfmpeg(input, output, uri, my_c)
			}
			m3u8 := "#EXTM3U\r\n#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=200000\r\nhttp://" +
				"42.51.201.196" + uri
			c.Ctx.WriteString(m3u8)
			return
		}
	}

	//TODO Return 404
	fmt.Println(uri)
}

func startFfmpeg(input string, output string, uri string, c chan string) {
	//time.Sleep(time.Second * 100)

	cmd := exec.Command("/bin/ffmpeg", "-rtsp_transport", "tcp", "-i",
		input, "-vcodec", "copy", "-acodec", "copy", "-f", "flv", output)

	stderr, err := cmd.StderrPipe()
	cmd.Start()
	content, err := ioutil.ReadAll(stderr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(content))

	c <- uri
}

func checkFfmpeg(m map[string]bool, c chan string) {
	//func checkFfmpeg(this *BeeferController) {
	for {
		select {
		case msg := <-c:
			m[msg] = false
			fmt.Println("received", msg)
		}
	}
}

func main() {
	var beeStruct BeeferController
	go checkFfmpeg(my_m, my_c)
	beego.Router("/hls/*", &beeStruct)
	beego.Run(":8000")
}
