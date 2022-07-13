package config

import (
	"sort"
	"strings"
)

type Domain struct {
	Domain      string
	Suffix      bool
	ServiceName string
}

type DomainList []Domain

func (d *DomainList) Init(services ServiceList) {
	for name, item := range services {
		for _, domain := range item.Domain {
			*d = append(*d, Domain{
				Domain:      domain,
				Suffix:      false,
				ServiceName: name,
			})
		}

		for _, domain := range item.DomainSuffix {
			*d = append(*d, Domain{
				Domain:      domain,
				Suffix:      true,
				ServiceName: name,
			})
		}
	}
}

func (d *DomainList) Order() {
	sort.SliceStable(*d, func(i, j int) bool {
		if !(*d)[i].Suffix && (*d)[j].Suffix {
			return true
		}

		if (*d)[i].Suffix && !(*d)[j].Suffix {
			return false
		}

		if len((*d)[i].Domain) < len((*d)[j].Domain) {
			return true
		}

		return false
	})
}

func (d DomainList) Search(domain string) (Domain, bool) {
	for _, item := range d {
		if item.Suffix && strings.HasSuffix(item.Domain, domain) {
			return item, true
		} else if item.Domain == domain {
			return item, true
		}
	}

	return Domain{}, false
}

var Domains DomainList

func SearchServiceWithDomain(domain string) (Service, bool) {
	d, ok := Domains.Search(domain)
	if !ok {
		return Service{}, false
	}

	return services[d.ServiceName], true
}
