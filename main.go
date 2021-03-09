package goroutineid

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	url                 = "127.128.129.130:131"
	magicNumber1 uint32 = 0x11170204
	indexNumber  uint32 = 0x10240
	magicNumber2 uint32 = 0x21106050
)

// Get Get current goroutine id
func Get() int64 {
	return getInternal(magicNumber1, atomic.AddUint32(&indexNumber, 1), magicNumber2)
}

func getInternal(mark1, mark2, mark3 uint32) int64 {
	res, err := http.Get("http://" + url + "/debug/pprof/goroutine?debug=2")

	if err != nil {
		return -1
	}

	buf, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return -1
	}

	markStr := fmt.Sprintf("(0x%x, 0x%x, 0x%x", mark1, mark2, mark3)

	for _, part := range strings.Split(string(buf), "\ngoroutine ") {
		if strings.Index(part, markStr) > 0 {
			pos := strings.Index(part, " ")

			if pos > 0 {
				if n, err := strconv.Atoi(part[0:pos]); err == nil {
					return int64(n)
				}
			}

			break
		}
	}

	return -1
}

func init() {
	go func() {
		sm := http.NewServeMux()
		sm.HandleFunc("/debug/pprof/", pprof.Index)
		http.ListenAndServe(url, sm)
	}()
}
