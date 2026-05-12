package forward

import (
	"net/http"
	"net/url"
)

var _ http.CookieJar = (*CookieJar)(nil)

type CookieJar struct {
	cookies map[string]*http.Cookie
}

func newCookieJar() *CookieJar {
	return &CookieJar{
		cookies: make(map[string]*http.Cookie),
	}
}

// SetCookies sets a list of cookies for the given request.
func (c *CookieJar) SetCookies(_ *url.URL, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		if cookie.Value == "deleted" {
			continue
		}

		c.cookies[cookie.Name] = cookie
	}
}

// GetCookies returns a list of cookies for the given request.
func (c *CookieJar) Cookies(_ *url.URL) []*http.Cookie {
	var ret []*http.Cookie
	for _, cookie := range c.cookies {
		ret = append(ret, cookie)
	}

	return ret
}
