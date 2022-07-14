package config

import (
	"strings"

	"gopkg.in/ini.v1"
)

type ServerUpgradeCommon struct {
	Address string `ini:"address"`
	Port    string `ini:"port"`
	SSL     bool   `ini:"ssl"`
	SSLCrt  string `ini:"ssl_crt,omitempty"`
	SSLKey  string `ini:"ssl_key,omitempty"`
	CheckMS uint32 `ini:"check_ms,omitempty"`
}

const (
	ServerWebSocketMode302PrefixURI = "/ws302/"

	ServerWebSocketModeBlock = "block"
	ServerWebSocketModeProxy = "proxy"
	ServerWebSocketMode302   = "302"
)

func ParseUpgrade(configFile string) (ServerUpgradeCommon, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowShadows: true,
	}, configFile)
	if err != nil {
		return globalServerUpgradeCommon, err
	}

	// config common
	commonSection, err := cfg.GetSection("common")
	if err != nil {
		return globalServerUpgradeCommon, err
	}

	if err := commonSection.MapTo(&globalServerUpgradeCommon); err != nil {
		return globalServerUpgradeCommon, err
	}

	// challenge port
	port := Port{}
	if err := port.Set(globalServerUpgradeCommon.Port); err != nil {
		return globalServerUpgradeCommon, err
	}
	globalServerUpgradeCommon.Port = port.Get()

	// config service
	services = make(ServiceList)

	for _, section := range cfg.Sections() {
		name := section.Name()
		if !strings.HasPrefix(name, "service.") {
			continue
		}

		service := Service{}
		if err := section.MapTo(&service); err != nil {
			return ServerUpgradeCommon{}, err
		}

		port := Port{}
		if err := port.Set(service.DestPort); err != nil {
			return globalServerUpgradeCommon, err
		}
		service.DestPort = port.Get()

		if service.CheckMS == 0 {
			service.CheckMS = 2000 // default check 2000 ms
		}

		services.Add(name, service)
	}

	Domains = []Domain{}
	Domains.Init(services)
	Domains.Order()

	return globalServerUpgradeCommon, nil
}

var globalServerUpgradeCommon ServerUpgradeCommon

func GetServerUpgradeCommon() ServerUpgradeCommon {
	return globalServerUpgradeCommon
}
