package dto

type ApplicationDTO struct {
	ID                    string `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	Type                  string `json:"type,omitempty"`
	SubType               string `json:"subType,omitempty"`
	IconUrl               string `json:"iconUrl,omitempty"`
	IsVisible             bool   `json:"isVisible"`
	IsNotificationEnabled bool   `json:"isNotificationEnabled"`
	Enabled               bool   `json:"enabled"`

	ConnectionSettings ConnectionSettingsDTO  `json:"connectionSettings"`
	TcpTunnelSettings  []TcpTunnelSettingsDTO `json:"tcpTunnelSettings"`
}

type ConnectionSettingsDTO struct {
	InternalAddress       string `json:"internalAddress"`
	Subdomain             string `json:"subdomain"`
	LuminateAddress       string `json:"luminateAddress,omitempty"`
	ExternalAddress       string `json:"externalAddress,omitempty"`
	CustomExternalAddress string `json:"customExternalAddress,omitempty"`
	CustomRootPath        string `json:"customRootPath,omitempty"`
	HealthUrl             string `json:"healthUrl,omitempty"`
	HealthMethod          string `json:"healthMethod,omitempty"`
	CustomSSLCertificate  string `json:"customSSLCertificate,omitempty"`
	WildcardPrivateKey    string `json:"wildcardPrivateKey,omitempty"`
}

type TcpTunnelSettingsDTO struct {
	Target      string        `json:"target"`
	Ports       []int         `json:"ports"`
	PortMapping []int         `json:"portMapping"`
	PortRanges  []interface{} `json:"portRanges"`
}

type HttpLinkTranslationSettingsDTO struct {
	IsDefaultContentRewriteRulesEnabled bool     `json:"isDefaultContentRewriteRulesEnabled"`
	IsDefaultHeaderRewriteRulesEnabled  bool     `json:"isDefaultHeaderRewriteRulesEnabled"`
	UseExternalAddressForHostAndSni     bool     `json:"useExternalAddressForHostAndSni"`
	LinkedApplications                  []string `json:"linkedApplications"`
}

type HttpRequestCustomizationSettings struct {
	HeaderCustomization map[string]string `json:"headerCustomization"`
}

type ApplicationPageDTO struct {
	First            bool             `json:"first"`
	Last             bool             `json:"last"`
	NumberOfElements int              `json:"numberOfElements"`
	Content          []ApplicationDTO `json:"content"`
	PageNumber       int              `json:"number"`
	PageSize         int              `json:"size"`
	TotalElements    int              `json:"totalElements"`
	TotalPages       int              `json:"totalPages"`
}
