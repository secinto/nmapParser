package parser

import (
	"github.com/antchfx/xmlquery"
	utils "secinto/checkfix_utils"
	"strconv"
)

func GetAllHostEntries(document *xmlquery.Node, key string) []Host {
	var hosts []Host
	var host Host

	hostRecords := utils.GetAllXMLNodesForKey(document, key)
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

func getGeneralInfoForHost(node *xmlquery.Node) Host {
	var host Host
	ipAddress := utils.GetValueForQuery(node, "//address/@addr")
	hostName := utils.GetValueForQuery(node, "//hostnames/hostname/@name")
	if ipAddress != "" {
		host.IP = ipAddress
		host.Name = hostName
	}
	return host
}

func getAllServiceForHost(node *xmlquery.Node, onlyOpen bool) []Service {
	var services []Service
	ports := utils.GetAllXMLNodesForKey(node, "ports/port")
	for _, port := range ports {
		var service Service
		portNumber := utils.GetValueForQuery(port, "//@portid")
		portProtocol := utils.GetValueForQuery(node, "//@protocol")
		if portNumber != "" && portProtocol != "" {
			service.Number, _ = strconv.Atoi(portNumber)
			service.Protocol = portProtocol
			service.State = utils.GetValueForQuery(port, "//state/@state")
			if onlyOpen && service.State != "open" {
				continue
			}
			service.Name = utils.GetValueForQuery(port, "//service/@name")
			service.Product = utils.GetValueForQuery(port, "//service/@product")
			service.Description = utils.GetValueForQuery(port, "//service/@extrainfo")
			service.OS = utils.GetValueForQuery(port, "//service/@ostype")

			services = append(services, service)
		}
	}
	return services
}
