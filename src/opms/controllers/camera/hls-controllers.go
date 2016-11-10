package camera

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

func getInt(str string) (int64, error) {
	intExp := myRegexp{regexp.MustCompile(`(?P<first>[\d]+)`)}
	mmap := intExp.FindStringSubmatchMap(str)
	if len(mmap) == 0 {
		return 0, errors.New("empty")
	}
	return strconv.ParseInt(mmap["first"], 10, 64)
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

func parseLiveStream(line string) (int64, int64) {
	var frames, bitrates int64
	stringMap := strings.Split(line, "=")                                   // 用 "=" 分隔每一列
	if len(stringMap) == 8 && strings.Compare(stringMap[0], "frame") == 0 { // 有8列，且第一列是"frame" 就是正常的直播流日志
		for i, v := range stringMap {
			if i == 1 { //  frame value
				if thisFrames, err := getInt(v); err == nil {
					frames = thisFrames
				}
			}
			if i == 6 { // bitrate value
				if thisBitrate, err := getInt(v); err == nil {
					bitrates = thisBitrate
				}
			}
		}
	}
	return frames, bitrates // 返回 总的视频帧数 码率
}

func catLogFile(reader *bufio.Reader) (int64, int64, int64) {
	var frames, bitrates int64 = 0, 0
	var count int64 = 0
	for {
		line, err := reader.ReadString('\r') // 读取 \r 分隔的行
		if err == io.EOF {
			break
		}
		count++
		frames, bitrates = parseLiveStream(strings.Replace(line, " ", "", -1))
	}
	return frames, bitrates, count
}

func startFfmpeg(input string, output string, logFile string) {
	if runtime.GOOS != "windows" {
		go func() {
			cmd := exec.Command("/bin/ffmpeg", "-re", "-rtsp_transport", "tcp", "-i", input,
				"-re", "-f", "lavfi", "-i", "anullsrc", "-c:v", "copy", "-c:a", "aac", "-map", "0:v", "-map", "1:a", "-f", "flv", output)

			if writeErrorLogFD, writeErrorLogErr := os.Create("log/" + logFile); writeErrorLogErr != nil {
				fmt.Println("Can not create log file : log.txt")
				return
			} else {
				defer writeErrorLogFD.Close()
				cmd.Stderr = writeErrorLogFD
			}

			readErrorLogFD, readErrorLogErr := os.OpenFile(("log/" + logFile), os.O_RDONLY, 0660)
			if readErrorLogErr != nil {
				fmt.Println("Can not create log file : log.txt")
				return
			}
			defer readErrorLogFD.Close()
			reader := bufio.NewReader(readErrorLogFD)

			if cmd.Start() != nil { // 启动进程
				fmt.Println("Can not start cmd.")
				return
			}

			var lastFrames, lastBitrates int64 = 0, 0
			var continuedNoLogTimes int = 0
			for { // 检测ffmpeg拉流是否正常
				if frames, bitrates, count := catLogFile(reader); count != 0 && frames > lastFrames {
					lastFrames = frames
					lastBitrates = bitrates
					continuedNoLogTimes = 0
				} else {
					continuedNoLogTimes++ // 读不到新的log 或者视频没有数据
				}
				fmt.Println(lastFrames, lastBitrates, continuedNoLogTimes)
				time.Sleep(time.Second * 1)

				if continuedNoLogTimes >= 10 { // 10秒都没有数据，或没有视频流，退出
					fmt.Println("Log not change, going to kill cmd.")
					cmd.Process.Kill()
					break
				}
			}

			cmd.Wait()
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
