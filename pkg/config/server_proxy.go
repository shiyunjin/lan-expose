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

type ServerProxyWebSocket struct {
	Mode302Domain string `ini:"mode_302_domain,omitempty"`
}

func ParseProxy(configFile string) (ServerProxyCommon, error) {
	cfg, err := ini.Load(configFile)
	if err != nil {
		return globalServerProxyCommon, err
	}

	// config common
	commonSection, err := cfg.GetSection("common")
	if err != nil {
		return globalServerProxyCommon, err
	}

	if err := commonSection.MapTo(&globalServerProxyCommon); err != nil {
		return globalServerProxyCommon, err
	}

	// config websocket
	webSocketSection, err := cfg.GetSection("websocket")
	if err != nil {
		return globalServerProxyCommon, err
	}

	if err := webSocketSection.MapTo(&globalServerProxyWebSocket); err != nil {
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

var globalServerProxyWebSocket ServerProxyWebSocket

func GetServerProxyWebSocket() ServerProxyWebSocket {
	return globalServerProxyWebSocket
}
