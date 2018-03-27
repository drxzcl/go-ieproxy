// +build !windows

package ieproxy

import (
	"net/url"
)

func (apc *AutomaticProxyConf) findProxyForURL(u *url.URL) (*url.URL, error) {
	return nil, nil
}
