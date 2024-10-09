package parser

const VERSION = "0.2.1"

type NmapParser struct {
	options *Options
}

type Config struct {
	ProjectsPath string `yaml:"projects_path"`
	PortsXMLFile string `yaml:"ports_xml,omitempty"`
	ServicesFile string `yaml:"services_file,omitempty"`
	HostMapping  string `yaml:"host_mapping,omitempty"`
}

type Project struct {
	Name string `yaml:"name"`
}

type Host struct {
	IP              string   `json:"ip"`
	Name            string   `json:"hostname"`
	Ports           []Port   `json:"ports"`
	AssociatedNames []string `json:"associatedNames"`
}

type Port struct {
	Number   string `json:"port"`
	Protocol string `json:"protocol"`
}

type Service struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service"`
}
