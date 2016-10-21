package main

//import "os"
import "os/exec"
import "fmt"
import "net/http"
import "bufio"
import "time"
import "regexp"

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

func getIPAddress(str string) (ip string, err bool) {
	var myExp = myRegexp{regexp.MustCompile(`var url='http://(?P<first>\d+)\.(?P<second>\d+).(?P<third>\d+).(?P<fourth>\d+):`)}
	mmap := myExp.FindStringSubmatchMap(str)
	ip = mmap["first"] + "." + mmap["second"] + "." + mmap["third"] + "." + mmap["fourth"]
	return ip, len(mmap) == 4
}

func execCommand() {
	cmd := exec.Command("date")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	cmd.StderrPipe()

	time.Sleep(time.Second)
}

func getDDNS(dns string) (ip string, ret bool) {
	resp, err := http.Get("http://www.hik-online.com/" + dns)
	if err != nil {
		return ip, false
	}

	defer resp.Body.Close()
	br := bufio.NewReader(resp.Body)
	line, isPrefix, err := br.ReadLine()
	for i := 1; err == nil && !isPrefix; i++ {
		if i == 113 {
			s := string(line)
			ip, ret = getIPAddress(s)
			if ret == true {
				return ip, true
			}
			return ip, false
		}
		line, isPrefix, err = br.ReadLine()
	}

	return ip, false
}

func main() {
	//execCommand()
	if ip, err := getDDNS("tnwl123456"); err == false {
		fmt.Println("HTTP Get Error")
	} else {
		fmt.Println(ip)
	}
}
