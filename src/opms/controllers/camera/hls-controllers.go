package camera

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"opms/controllers"
	. "opms/models/camera"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
)

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

func parseURI(uri string) (int64, string, string, string, bool) {
	myExp := myRegexp{regexp.MustCompile(`/hls/(?P<first>[\w]+|[\d]+)_(?P<second>[\w]+|[\d])_(?P<third>[\d]+).m3u8`)}
	mmap := myExp.FindStringSubmatchMap(uri)
	id, _ := strconv.ParseInt(mmap["third"], 10, 64)
	return id, mmap["first"], mmap["second"], ("/hls/" + mmap["first"] + "_" + mmap["second"] + "_" + mmap["third"]), len(mmap) != 3
}

func getIPAddress(str string) (string, bool) {
	myExp := myRegexp{regexp.MustCompile(`var url='http://(?P<first>\d+)\.(?P<second>\d+).(?P<third>\d+).(?P<fourth>\d+):`)}
	mmap := myExp.FindStringSubmatchMap(str)
	ip := mmap["first"] + "." + mmap["second"] + "." + mmap["third"] + "." + mmap["fourth"]
	return ip, len(mmap) == 4
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

func isFfmpegStartUp(dns string, channel string) error {
	if runtime.GOOS != "windows" {
		mycmd, _ := os.Getwd()
		cmd := exec.Command(mycmd+"/isFfmpegStartUp.sh", dns, channel)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("will not start")
			return errors.New("will not start")
		}
		if len(out) > 0 {
			fmt.Println(string(out))
			return nil
		}
		fmt.Println("will not start")
		return errors.New("will not start")
	}
	return errors.New("os system not support")
}

func startFfmpeg(input string, output string) {
	if runtime.GOOS != "windows" {
		cmd := exec.Command("/bin/ffmpeg", "-rtsp_transport", "tcp", "-i",
			input, "-vcodec", "copy", "-acodec", "copy", "-f", "flv", output)

		go func() {
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("end")
		}()
	}
}

type HLSController struct {
	controllers.BaseController
}

func (this *HLSController) Get() {
	uri := this.Ctx.Request.URL.String()

	if id, dns, channel, rtmpURI, err := parseURI(uri); err != true {
		if ip, err := getDdnsIP(dns); err == true {
			// 判断id是否存在
			cam, err := GetDeparts(int64(id))
			if err != nil {
				this.Abort("404")
				return
			}
			if err := isFfmpegStartUp(dns, channel); err == nil {
				input := "rtsp://" + cam.Users + ":" + cam.Pass + "@" + ip + ":1554/mpeg4/" + channel + "/sub/av_stream"
				output := "rtmp://42.51.201.196" + rtmpURI
				startFfmpeg(input, output)

				fmt.Println(input, "--", output)
				m3u8 := "#EXTM3U\r\n#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=200000\r\nhttp://" +
					"42.51.201.196" + rtmpURI + ".m3u8"
				this.Ctx.WriteString(m3u8)
				return
			}
		}
	}
	this.Abort("404")
}
