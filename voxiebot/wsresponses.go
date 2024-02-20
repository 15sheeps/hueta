package main

type WebsocketResponse struct {
	Type            string      `yaml:"type"`
	MatchID         string      `yaml:"matchId"`
	Namespace       string      `yaml:"namespace"`
	PodName         string      `yaml:"podName"`
	Region          string      `yaml:"region"`
	ImageVersion    string      `yaml:"imageVersion"`
	IP              string      `yaml:"ip"`
	Port            int         `yaml:"port"`
	Ports           Ports       `yaml:"ports"`
	Protocol        string      `yaml:"protocol"`
	ServerVersion   string      `yaml:"serverVersion"`
	Provider        string      `yaml:"provider"`
	Status          string      `yaml:"status"`
	Deployment      string      `yaml:"deployment"`
	CustomAttribute string      `yaml:"customAttribute"`
	IsOK            bool        `yaml:"isOK"`
	Message         interface{} `yaml:"message"`
}
type Ports struct {
}