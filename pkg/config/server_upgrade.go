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
}

func ParseUpgrade(configFile string) (ServerUpgradeCommon, error) {
	cfg, err := ini.Load(configFile)
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

		port := Port{}
		if err := port.Set(section.Key("dest_port").Value()); err != nil {
			return globalServerUpgradeCommon, err
		}

		services.Add(name, Service{
			Domain:       section.Key("domain").ValueWithShadows(),
			DomainSuffix: section.Key("domain_suffix").ValueWithShadows(),
			DestDomain:   section.Key("dest_domain").Value(),
			DestPort:     port.Get(),
		})
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
