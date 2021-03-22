package net

import (
	"fmt"
	"math"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// IPv4 Segment
	v4Seg   = "(?:[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])"
	v4Str   = fmt.Sprintf("(%v[.]){3}%v", v4Seg, v4Seg)
	IPv4Reg = regexp.MustCompile("^" + v4Str + "$")

	// IPv6 Segment
	v6Seg   = "(?:[0-9a-fA-F]{1,4})"
	IPv6Reg = regexp.MustCompile(
		fmt.Sprintf(`^(`+
			`(?:%v:){7}(?:%v|:)|`+
			`(?:%v:){6}(?:%v|:%v|:)|`+
			`(?:%v:){5}(?::%v|(:%v){1,2}|:)|`+
			`(?:%v:){4}(?:(:%v){0,1}:%v|(:%v){1,3}|:)|`+
			`(?:%v:){3}(?:(:%v){0,2}:%v|(:%v){1,4}|:)|`+
			`(?:%v:){2}(?:(:%v){0,3}:%v|(:%v){1,5}|:)|`+
			`(?:%v:){1}(?:(:%v){0,4}:%v|(:%v){1,6}|:)|`+
			`(?::((?::%v){0,5}:%v|(?::%v){1,7}|:))`,
			v6Seg, v6Seg,
			v6Seg, v4Str, v6Seg,
			v6Seg, v4Str, v6Seg,
			v6Seg, v6Seg, v4Str, v6Seg,
			v6Seg, v6Seg, v4Str, v6Seg,
			v6Seg, v6Seg, v4Str, v6Seg,
			v6Seg, v6Seg, v4Str, v6Seg,
			v6Seg, v4Str, v6Seg,
		) + `)(%[0-9a-zA-Z]{1,})?$`)

	// 全局 socket 结束 WaitGroup
	SocketWaitGroup sync.WaitGroup
	// debug socketWaitGroup 状态
	debugSocketWaitGroup bool
)

func IsIPv4(s string) bool {
	return IPv4Reg.MatchString(s)
}

func IsIPv6(s string) bool {
	return IPv6Reg.MatchString(s)
}

func IsIP(s string) int {
	if IsIPv4(s) {
		return 4
	} else if IsIPv6(s) {
		return 6
	} else {
		return 0
	}
}

func IsLegalPort(port int) bool {
	if port < 0 {
		return false
	}
	return int(math.Abs(float64(port))) == (int(math.Abs(float64(port)))>>0) && port <= 0xFFFF
}

func GetSocketWaitGroup(tags ...string) *sync.WaitGroup {
	// 如果 打开 debug waitGroup
	if GetDebugSocketWaitGroup() {
		if len(tags) > 0 {
			_, file, line, ok := runtime.Caller(1)
			var code string
			if ok {
				code = file + ":" + strconv.Itoa(line)
			} else {
				code = ""
			}
			fmt.Printf("\033[47;30m%s\033[0m %s\n", code, strings.Join(tags, " "))
		}
	}
	return &SocketWaitGroup
}

func ExitAfterSocketEnd() {
	GetSocketWaitGroup().Wait()
}

func WaitAfterSocketEnd() {
	for {
		GetSocketWaitGroup().Wait()
		time.Sleep(50 * time.Millisecond)
	}
}

// 激活 debugSocketWaitGroup
func EnableDebugSocketWaitGroup() {
	debugSocketWaitGroup = true
}

// 禁用 debugSocketWaitGroup
func DisableDebugSocketWaitGroup() {
	debugSocketWaitGroup = false
}

// 获取 debugSocketWaitGroup 状态
func GetDebugSocketWaitGroup() bool {
	return debugSocketWaitGroup
}
