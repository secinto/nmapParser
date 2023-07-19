package parser

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

var (
	log           = NewLogger()
	appConfig     Config
	project       Project
	excludedPorts = []int{21, 22, 25, 143, 110, 143, 465, 587, 993, 995}
	includedPorts = []int{80, 81, 280, 443, 591, 593, 445, 457, 832, 981, 1311, 1241,
		1342, 1433, 1434, 1443, 1521, 1944, 2301, 2080, 2443, 2480, 3000, 3080, 3128,
		3306, 3443, 4000, 4001, 4002, 4080, 4100, 4443, 4567, 5000, 5080, 5104, 5200,
		5432, 5443, 5800, 5801, 5802, 5800, 6080, 6346, 6347, 6443, 7001, 7002, 7021,
		7023, 7025, 7080, 7443, 7777, 8000, 8008, 8042, 8080, 8081, 8082, 8180, 8222,
		8280, 8281, 8333, 8443, 8530, 8531, 8888, 8887, 9000, 9080, 9090, 9443, 10443,
		11443, 12443, 13443, 14443, 15443, 16080, 30821}
)

/*
--------------------------------------------------------------------------------

	Initialization functions for the application

-------------------------------------------------------------------------------
*/
func (p *NmapParser) initialize(configLocation string) {
	appConfig = loadConfigFrom(configLocation)
	if !strings.HasSuffix(appConfig.S2SPath, "/") {
		appConfig.S2SPath = appConfig.S2SPath + "/"
	}
	p.options.BaseFolder = appConfig.S2SPath + p.options.Project
	if !strings.HasSuffix(p.options.BaseFolder, "/") {
		p.options.BaseFolder = p.options.BaseFolder + "/"
	}
	appConfig.PortsXMLFile = strings.Replace(appConfig.PortsXMLFile, "{project_name}", p.options.Project, -1)
	//appConfig.DPUXOutput = strings.Replace(appConfig.DPUXOutput, "{project_name}", p.options.Project, -1)

	project = Project{
		Name: p.options.Project,
	}

}

func loadConfigFrom(location string) Config {
	var config Config
	var yamlFile []byte
	var err error

	yamlFile, err = os.ReadFile(location)
	if err != nil {
		yamlFile, err = os.ReadFile(defaultSettingsLocation)
		if err != nil {
			log.Fatalf("yamlFile.Get err   #%v ", err)
		}
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	if &config == nil {
		config = Config{
			S2SPath:          "S://",
			DPUXCleanXMLFile: "dpux_clean.xml",
			PortsXMLFile:     "ports.{project_name}.output.xml",
			HostMapping:      "dpux_host_to_ip.json",
		}
	}
	return config
}

func NewParser(options *Options) (*NmapParser, error) {
	parser := &NmapParser{options: options}
	parser.initialize(options.SettingsFile)
	return parser, nil
}

func (p *NmapParser) Parse() error {
	if p.options.Project != "" {
		if p.options.ServiceJSON {
			log.Infof("Parsing NMAP output for project %s", p.options.Project)
			allHostRecords := p.parseDPUX()
			WriteHostsToJSONFile(p.options.BaseFolder+"findings/service_test.json", allHostRecords)
		}
		if p.options.HTTPDomains {
			log.Infof("Creating host name mapping including possible open HTTP ports for project %s", p.options.Project)
			hostEntries := p.generateHostPortCombinations(false)
			WriteToFile(p.options.BaseFolder+"domains_with_http_ports.txt", hostEntries)
		}
		if p.options.All {
			log.Infof("Creating host name mapping including all ports for project %s", p.options.Project)
			hostEntries := p.generateHostPortCombinations(true)
			WriteToFile(p.options.BaseFolder+"domains_with_ports.txt", hostEntries)
		}
	}
	return nil
}

/*
--------------------------------------------------------------------------------

	Public functions of the application

-------------------------------------------------------------------------------
*/
func (p *NmapParser) parseDPUX() []Host {
	input := GetXMLDocumentFromFile(p.options.BaseFolder + "nmap/" + appConfig.DPUXCleanXMLFile)
	allHostRecords := GetAllHostEntries(input, "host")
	if log.Level == logrus.InfoLevel {
		if len(allHostRecords) >= 1 {
			// Fine, we found at least one.
			for _, hostNode := range allHostRecords {
				if hostNode.Name != "" {
					log.Infof("Found host with IP %s and name %s with %d services running", hostNode.IP, hostNode.Name, len(hostNode.Services))
				} else {
					log.Infof("Found host with IP %s with %d services running", hostNode.IP, len(hostNode.Services))
				}
			}
		}
	}
	return allHostRecords
}

func (p *NmapParser) parsePorts() []Host {
	input := GetXMLDocumentFromFile(p.options.BaseFolder + "recon/" + appConfig.PortsXMLFile)
	allHostRecords := GetAllHostEntries(input, "host")
	if log.Level == logrus.InfoLevel {
		if len(allHostRecords) >= 1 {
			// Fine, we found at least one.
			for _, hostNode := range allHostRecords {
				if hostNode.Name != "" {
					log.Infof("Found host with IP %s and name %s with %d services running", hostNode.IP, hostNode.Name, len(hostNode.Services))
				} else {
					log.Infof("Found host with IP %s with %d services running", hostNode.IP, len(hostNode.Services))
				}
			}
		}
	}
	return allHostRecords
}

func (p *NmapParser) generateHostPortCombinations(generateAll bool) string {
	input := GetJSONDocumentFromFile(p.options.BaseFolder + "recon/" + appConfig.HostMapping)
	allIPEntries := GetAllRecordsForKey(input, "ip")
	allIPHosts := make(map[string]Host)

	if len(allIPEntries) >= 1 {
		// Fine, we found at least one.
		for _, host := range allIPEntries {
			hostIPEntries := getValuesFromNode(host, "ip")
			hostNameEntries := getValuesFromNode(host, "host")
			if len(hostIPEntries) >= 1 {
				for _, hostEntry := range hostIPEntries {
					if existingHost, ok := allIPHosts[hostEntry]; !ok {
						// First entry should be created not from dpux but rather from dns resolution of nmap possibly
						newHost := Host{
							IP:              hostEntry,
							Name:            hostNameEntries[0],
							Services:        nil,
							AssociatedNames: []string{},
						}
						allIPHosts[hostEntry] = newHost
					} else {
						existingHost.AssociatedNames = AppendIfMissing(existingHost.AssociatedNames, hostNameEntries[0])
						allIPHosts[hostEntry] = existingHost
					}
				}
			}
		}
	}
	allPortsRecords := p.parsePorts()

	for _, ports := range allPortsRecords {
		if existingHost, ok := allIPHosts[ports.IP]; ok {
			existingHost.Services = ports.Services
			allIPHosts[ports.IP] = existingHost
		}
	}
	var counter = 0
	var hostEntries strings.Builder
	for _, host := range allIPHosts {
		for _, service := range host.Services {
			if checkIfPortIsContained(service.Number, includedPorts) || generateAll {
				hostEntries.WriteString(host.Name + ":" + strconv.Itoa(service.Number) + "\n")
				counter++
				for _, additionalHost := range host.AssociatedNames {
					hostEntries.WriteString(additionalHost + ":" + strconv.Itoa(service.Number) + "\n")
					counter++
				}
			}
		}
	}

	log.Infof("Created %d mappings from %d initial domain names", counter, len(allIPEntries))
	return hostEntries.String()
}
