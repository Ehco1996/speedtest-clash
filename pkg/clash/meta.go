package clash

import (
	"fmt"
	"net/url"

	"github.com/Dreamacro/clash/constant"
)

func URLToMetadata(u *url.URL) (*constant.Metadata, error) {
	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			return nil, fmt.Errorf("%s scheme not Support", u.String())

		}
	}
	meta := &constant.Metadata{
		AddrType: constant.AtypDomainName,
		Host:     u.Hostname(),
		DstIP:    nil,
		DstPort:  port,
	}
	return meta, nil
}
