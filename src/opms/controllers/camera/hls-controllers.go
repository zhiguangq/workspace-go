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
	"strings"
	"time"

	"github.com/astaxie/beego"
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
		cmd := exec.Command("ps", "-ef")
		out, _ := cmd.CombinedOutput()
		lines := strings.Split(string(out), "\n")
		for _, s := range lines {
			if strings.Contains(s, dns) && strings.Contains(s, channel) {
				//fmt.Println("Has been start : " + s)
				return errors.New("Has been started.")
			}
		}
		return nil
	}
	return errors.New("os system not support")
}

func startFfmpeg(input string, output string, logFile string) {
	if runtime.GOOS != "windows" {
		go func() {
			cmd := exec.Command("/bin/ffmpeg", "-rtsp_transport", "tcp", "-i", input,
				"-f", "lavfi", "-i", "anullsrc", "-c:v", "copy", "-c:a", "aac", "-map", "0:v", "-map", "1:a", "-f", "flv", output)
			errLog, err := os.Create("log/" + logFile)
			if err != nil {
				fmt.Println("Can not create log file : " + logFile)
				return
			}
			defer errLog.Close()
			cmd.Stderr = errLog

			err = cmd.Start()
			if err != nil {
				fmt.Println("Can not start cmd.")
				return
			}

			var log_size int64 = 0
			for {
				time.Sleep(time.Second * 10)
				info, _ := errLog.Stat()
				if log_size == info.Size() {
					fmt.Println("Log not change, going to kill cmd.")
					cmd.Process.Kill()
					break
				}
				log_size = info.Size()
			}
			cmd.Wait()
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
		if err := isFfmpegStartUp(dns, channel); err == nil {
			if ip, err := getDdnsIP(dns); err == true {
				// 判断id是否存在
				cam, err := GetDeparts(int64(id))
				if err != nil {
					this.Abort("404")
					return
				}

				input := "rtsp://" + cam.Users + ":" + cam.Pass + "@" + ip + ":1554/mpeg4/" + channel + "/sub/av_stream"
				output := "rtmp://" + beego.AppConfig.String("hlsServer") + ":" + beego.AppConfig.String("hlsRtmpPort") + rtmpURI
				startFfmpeg(input, output, (dns + "_" + channel))
				fmt.Println(input, "--", output)
			}
		}
		m3u8 := "#EXTM3U\r\n#EXT-X-STREAM-INF:PROGRAM-ID=1, BANDWIDTH=200000\r\nhttp://" +
			beego.AppConfig.String("hlsServer") + ":" + beego.AppConfig.String("hlsServerPort") + rtmpURI + ".m3u8"
		this.Ctx.WriteString(m3u8)
		return

	}
	this.Abort("404")
}
