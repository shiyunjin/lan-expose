package config

type ServiceList map[string]Service

type Service struct {
	Domain       []string `ini:"domain,omitempty,allowshadow"`
	DomainSuffix []string `ini:"domain_suffix,omitempty,allowshadow"`
	DestDomain   string   `ini:"dest_domain"`
	DestPort     string   `ini:"dest_port"`

	WebSocketMode          string `ini:"websocket_mode,omitempty"`
	WebSocketMode302Domain string `ini:"websocket_mode_302_domain,omitempty"`

	CheckMS uint32 `ini:"check_ms,omitempty"`
}

func (s *ServiceList) Add(name string, service Service) {
	(*s)[name] = service
}

var services ServiceList
