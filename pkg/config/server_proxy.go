package config

import (
	"net/url"
	"strings"

	"gopkg.in/ini.v1"
)

type ServerProxyCommon struct {
	Address string `ini:"address"`
	Port    string `ini:"port"`
	SSLCrt  string `ini:"ssl_crt"`
	SSLKey  string `ini:"ssl_key"`
}

func ParseProxy(configFile string) (ServerProxyCommon, error) {
	cfg, err := ini.Load(configFile)
	if err != nil {
		return globalServerProxyCommon, err
	}

	// config common
	if err := cfg.MapTo(&globalServerProxyCommon); err != nil {
		return globalServerProxyCommon, err
	}

	// challenge port
	port := Port{}
	if err := port.Set(globalServerProxyCommon.Port); err != nil {
		return globalServerProxyCommon, err
	}
	globalServerProxyCommon.Port = port.Get()

	// config service
	proxys = ProxyList{
		list: map[string]*url.URL{},
	}

	for _, section := range cfg.Sections() {
		name := section.Name()
		if !strings.HasPrefix(name, "proxy.") {
			continue
		}

		if err := proxys.Set(section.Key("domain").Value(), section.Key("target").Value()); err != nil {
			return ServerProxyCommon{}, err
		}
	}

	return globalServerProxyCommon, nil
}

var globalServerProxyCommon ServerProxyCommon

func GetServerProxyCommon() ServerProxyCommon {
	return globalServerProxyCommon
}
