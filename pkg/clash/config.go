package clash

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	cc "github.com/Dreamacro/clash/config"
	"github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
)

func readConfigFromFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("configuration file %s is empty", path)
	}
	return data, err
}

func readConfigFromHttp(url string) ([]byte, error) {
	var httpc = &http.Client{Timeout: 10 * time.Second}
	resp, err := httpc.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func LoadConfig(path string) (*cc.Config, error) {
	var buf []byte
	var err error

	if strings.Contains(path, "http") {
		buf, err = readConfigFromHttp(path)
	} else {
		buf, err = readConfigFromFile(path)
	}
	if err != nil {
		return nil, err
	}
	cfg, err := executor.ParseWithBytes(buf)
	if err != nil {
		return nil, err
	}

	// only keep direct and other proxy node
	for k, p := range cfg.Proxies {
		if p.Type() == constant.Reject {
			delete(cfg.Proxies, k)
		}
		if p.Type() == constant.Relay {
			delete(cfg.Proxies, k)
		}
		if p.Type() == constant.Selector {
			delete(cfg.Proxies, k)
		}
		if p.Type() == constant.Fallback {
			delete(cfg.Proxies, k)
		}
		if p.Type() == constant.URLTest {
			delete(cfg.Proxies, k)
		}
		if p.Type() == constant.LoadBalance {
			delete(cfg.Proxies, k)
		}
	}
	return cfg, nil
}
