package parser

import (
	"github.com/antchfx/xmlquery"
	"os"
	"strconv"
	"strings"
)

func GetXMLDocumentFromFile(filename string) *xmlquery.Node {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Reading JSON input file failed: %s %s", err.Error(), filename)
	}
	xmlReader := strings.NewReader(string(data))
	input, err := xmlquery.Parse(xmlReader)
	if err != nil {
		log.Fatalf("Reading JSON input file failed: %s %s", err.Error(), filename)
	}

	return input
}

func GetAllHostEntries(document *xmlquery.Node, key string) []Host {
	var hosts []Host
	var host Host

	hostRecords := getAllNodesForKey(document, key)
	for _, hostElement := range hostRecords {
		if hostElement != nil {

			host = getGeneralInfoForHost(hostElement)

			if host.IP != "" {
				services := getAllServiceForHost(hostElement, true)
				if services != nil {
					host.Services = services
				}
				hosts = append(hosts, host)
			}
		}
	}
	return hosts
}

func getAllNodesForKey(document *xmlquery.Node, key string) []*xmlquery.Node {
	entries, error := xmlquery.QueryAll(document, "//"+key)

	if error != nil {
		log.Errorf("Querying XML error #%v ", error)
	}

	return entries
}

func getGeneralInfoForHost(node *xmlquery.Node) Host {
	var host Host
	ipAddress := getValueForQuery(node, "//address/@addr")
	hostName := getValueForQuery(node, "//hostnames/@name")
	if ipAddress != "" {
		host.IP = ipAddress
		host.Name = hostName
	}
	return host
}

func getAllServiceForHost(node *xmlquery.Node, onlyOpen bool) []Service {
	var services []Service
	ports := getAllNodesForKey(node, "ports/port")
	for _, port := range ports {
		var service Service
		portNumber := getValueForQuery(port, "//@portid")
		portProtocol := getValueForQuery(node, "//@protocol")
		if portNumber != "" && portProtocol != "" {
			service.Number, _ = strconv.Atoi(portNumber)
			service.Protocol = portProtocol
			service.State = getValueForQuery(port, "//state/@state")
			if onlyOpen && service.State != "open" {
				continue
			}
			service.Name = getValueForQuery(port, "//service/@name")
			service.Product = getValueForQuery(port, "//service/@product")
			service.Description = getValueForQuery(port, "//service/@extrainfo")
			service.OS = getValueForQuery(port, "//service/@ostype")

			services = append(services, service)
		}
	}
	return services
}

func getValueForQuery(node *xmlquery.Node, query string) string {
	element, error := xmlquery.Query(node, query)

	if error != nil {
		log.Errorf("Querying XML error #%v ", error)
	}
	if element != nil {
		if element.InnerText() != "" {
			return element.InnerText()
		}
	}
	return ""
}
