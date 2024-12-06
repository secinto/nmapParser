package parser

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"os"
	utils "secinto/checkfix_utils"
	"strings"
)

var (
	log           = utils.NewLogger()
	appConfig     Config
	project       Project
	excludedPorts = []int{21, 22, 25, 143, 110, 143, 465, 587, 993, 995}
	includedPorts = []int{80, 81, 280, 443, 591, 593, 445, 457, 832, 981, 1311, 1241,
		1342, 1433, 1434, 1443, 1521, 1944, 2301, 2080, 2443, 2480, 3000, 3080, 3128,
		3306, 3443, 4000, 4001, 4002, 4080, 4100, 4343, 4443, 4567, 5000, 5080, 5104, 5200,
		5432, 5443, 5800, 5801, 5802, 5800, 6080, 6346, 6347, 6443, 7001, 7002, 7021,
		7023, 7025, 7080, 7443, 7777, 8000, 8008, 8042, 8080, 8081, 8082, 8180, 8222,
		8280, 8281, 8333, 8443, 8530, 8531, 8843, 8880, 8888, 8887, 9000, 9080, 9090, 9443, 10443,
		11443, 12443, 13443, 14443, 15443, 16080, 22443, 33443, 30821}
)

/*
--------------------------------------------------------------------------------

	Initialization functions for the application

-------------------------------------------------------------------------------
*/
func (p *NmapParser) initialize(configLocation string) {
	appConfig = loadConfigFrom(configLocation)
	if !strings.HasSuffix(appConfig.ProjectsPath, "/") {
		appConfig.ProjectsPath = appConfig.ProjectsPath + "/"
	}
	p.options.BaseFolder = appConfig.ProjectsPath + p.options.Project
	if !strings.HasSuffix(p.options.BaseFolder, "/") {
		p.options.BaseFolder = p.options.BaseFolder + "/"
	}

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
			ProjectsPath: "/checkfix/projects",
			ServicesFile: "services.json",
			HostMapping:  "dpux_host_to_ip.json",
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
		if p.options.HTTPDomains {
			log.Infof("Creating host name mapping including possible open HTTP ports for project %s", p.options.Project)
			hostEntries := p.generateHostPortCombinations(false)
			utils.WriteToFile(p.options.BaseFolder+"domains_with_http_ports.txt", hostEntries)
		}
		if p.options.All {
			log.Infof("Creating host name mapping including all ports for project %s", p.options.Project)
			hostEntries := p.generateHostPortCombinations(true)
			utils.WriteToFile(p.options.BaseFolder+"domains_with_ports.txt", hostEntries)
		}
	}
	return nil
}

/*
--------------------------------------------------------------------------------

	Public functions of the application

-------------------------------------------------------------------------------
*/
func (p *NmapParser) getServicesJSON() []Service {
	var services []Service
	servicesFile := p.options.BaseFolder + "findings/" + appConfig.ServicesFile
	servicesFileString := utils.ReadFileToString(servicesFile)
	if err := json.Unmarshal([]byte(servicesFileString), &services); err != nil {
		log.Fatalf("Couldn't unmarshal file %s. Reason: %v", servicesFile, err)
	}
	return services
}

func (p *NmapParser) generateHostPortCombinations(generateAll bool) string {
	input := utils.GetJSONDocumentFromFile(p.options.BaseFolder+"recon/"+appConfig.HostMapping, true)
	allIPEntries := utils.GetAllJSONNodesForKey(input, "ip")
	allIPHosts := make(map[string]Host)

	if len(allIPEntries) >= 1 {
		// Fine, we found at least one.
		for _, host := range allIPEntries {
			hostIPEntries := utils.GetValuesFromNode(host, "ip")
			hostNameEntries := utils.GetValuesFromNode(host, "host")
			for _, hostEntry := range hostIPEntries {
				if existingHost, ok := allIPHosts[hostEntry]; !ok {
					// First entry should be created not from dpux but rather from dns resolution of nmap possibly
					newHost := Host{
						IP:              hostEntry,
						Name:            strings.ToLower(hostNameEntries[0]),
						Ports:           nil,
						AssociatedNames: []string{},
					}
					allIPHosts[hostEntry] = newHost
				} else {
					existingHost.AssociatedNames = utils.AppendIfMissing(existingHost.AssociatedNames, strings.ToLower(hostNameEntries[0]))
					allIPHosts[hostEntry] = existingHost
				}
			}
		}
	}

	servicesJSON := p.getServicesJSON()

	for _, service := range servicesJSON {
		if existingHost, ok := allIPHosts[service.IP]; ok {
			testPort := Port{
				Number:   service.Port,
				Protocol: service.Protocol,
			}
			if len(existingHost.Ports) == 0 || checkIfPortIsNotContained(testPort, existingHost.Ports) {
				existingHost.Ports = append(existingHost.Ports, testPort)
			}
			allIPHosts[service.IP] = existingHost
		} else {
			log.Errorf("No services found for host with IP %s", service.IP)
		}
	}
	var counter = 0
	var hostEntries []string
	for _, host := range allIPHosts {
		for _, service := range host.Ports {
			if checkIfPortIsContained(service.Number, includedPorts) || generateAll {
				hostEntries = utils.AppendIfMissing(hostEntries, host.Name+":"+service.Number)
				counter++
				for _, additionalHost := range host.AssociatedNames {
					hostEntries = utils.AppendIfMissing(hostEntries, additionalHost+":"+service.Number)
					counter++
				}
			}
		}
	}
	var allHosts []Host
	for _, host := range allIPHosts {
		allHosts = append(allHosts, host)
	}
	if generateAll {
		data, _ := json.MarshalIndent(allHosts, "", " ")
		utils.WriteToFile(p.options.BaseFolder+"findings/all_dpux.json", string(data))
	}

	log.Infof("Created %d mappings from %d initial domain names", counter, len(allIPEntries))
	return strings.Join(hostEntries, "\n")
}
