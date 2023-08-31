package parser

const VERSION = "0.1"

type NmapParser struct {
	options *Options
}

type Config struct {
	S2SPath      string `yaml:"s2s_path"`
	PortsXMLFile string `yaml:"ports_xml,omitempty"`
	HostMapping  string `yaml:"host_mapping,omitempty"`
}

type Project struct {
	Name string `yaml:"name"`
}

type Host struct {
	IP              string    `json:"ip"`
	Name            string    `json:"hostname""`
	Services        []Service `json:"services"`
	AssociatedNames []string  `json:"associatedNames"`
}

type Service struct {
	Number      int    `json:"port"`
	Protocol    string `json:"protocol"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Product     string `json:"product,omitempty"`
	Description string `json:"description,omitempty"`
	OS          string `json:"ostype,omitempty"`
}
