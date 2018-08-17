package tools

import (
	"bytes"
	"os/exec"
	"sync"
)

var (
	cache sync.Map
)

func IptoUrl(ip string) string {
	var url string
	var err error
	if val, ok := cache.Load(ip); ok {
		url, _ = val.(string)
	} else {
		url, err = iptourl(ip)
		if err != nil {
			url = ""
			return url
		}
		cache.Store(ip, url)
	}
	return url
}

func iptourl(ip string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("./IP2Url.py", ip)
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}
