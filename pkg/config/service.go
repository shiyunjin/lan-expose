package config

type ServiceList map[string]Service

type Service struct {
	Domain       []string
	DomainSuffix []string
	DestDomain   string
	DestPort     string

	WebSocketMode          string
	WebSocketMode302Domain string
}

func (s *ServiceList) Add(name string, service Service) {
	(*s)[name] = service
}

var services ServiceList
