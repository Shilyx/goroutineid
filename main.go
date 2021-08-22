package goroutineid

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	magicNumber1 uint32 = 0x11170204
	magicNumber2 uint32 = 0x21106050
)

var (
	indexNumber uint32 = 0x10240
)

type respWriter struct {
	header     http.Header
	statusCode int
	buf        []byte
}

func (resp *respWriter) Header() http.Header {
	return resp.header
}

func (resp *respWriter) Write(buf []byte) (int, error) {
	resp.buf = append(resp.buf, buf...)
	return len(buf), nil
}

func (resp *respWriter) WriteHeader(statusCode int) {
	resp.statusCode = statusCode
}

// Get Get current goroutine id
func Get() int64 {
	return getInternal(magicNumber1, atomic.AddUint32(&indexNumber, 1), magicNumber2)
}

func getInternal(mark1, mark2, mark3 uint32) int64 {
	resp := respWriter{
		header:     http.Header(make(map[string][]string)),
		statusCode: http.StatusOK,
		buf:        []byte{},
	}
	req := http.Request{
		Method: "GET",
		URL: &url.URL{
			Path:     "/debug/pprof/goroutine",
			RawQuery: "debug=2",
		},
		Header: map[string][]string{
			"Accept":          {"text/plain"},
			"Accept-Encoding": {"identity"},
			"User-Agent":      {"github.com/shilyx/goroutineid"},
		},
	}
	pprof.Index(&resp, &req)

	if resp.statusCode != http.StatusOK {
		return -1
	}

	markStr := fmt.Sprintf("(0x%x, 0x%x,", mark1, mark2)
	markStr64 := fmt.Sprintf("(0x%x,", uint64(mark2)*0x100000000+uint64(mark1))

	for _, part := range strings.Split(string(resp.buf), "goroutine ") {
		if strings.Index(part, markStr) > 0 || strings.Index(part, markStr64) > 0 {
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
