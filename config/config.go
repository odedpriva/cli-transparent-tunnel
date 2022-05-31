package config

type Config struct {
	SacSettings *SacSettings
}

type SacSettings struct {
	ClientID     string
	ClientSecret string
	TenantDomain string
}

func GetConfig() (*Config, error) {
	return &Config{
		SacSettings: &SacSettings{
			ClientID:     "22f0aa9ef0ac982ad34085ef54bd9440",
			ClientSecret: "654de200af286f41779941ad02b9e41859f120c58b7c28a233d42aa1a420ca14",
			TenantDomain: "symchatbotdemo.luminatesite.com",
		},
	}, nil
}
