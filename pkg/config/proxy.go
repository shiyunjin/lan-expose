package config

import (
	"errors"
	"net/url"
)

type Proxy struct {
	Domain string
	Target string
}

type ProxyList struct {
	list map[string]*url.URL
}

func (p *ProxyList) Set(domain, target string) error {
	if _, ok := (*p).list[domain]; ok {
		return errors.New("duplicate domain: " + domain)
	}

	u, err := url.Parse(target)
	if err != nil {
		return err
	}

	(*p).list[domain] = u

	return nil
}

func (p ProxyList) Get(domain string) (*url.URL, bool) {
	u, ok := p.list[domain]

	return u, ok
}

var proxys ProxyList

func SearchTargetWithDomain(domain string) (*url.URL, bool) {
	return proxys.Get(domain)
}
