package nmapParser

import (
	"bufio"
	"encoding/json"
	"golang.org/x/exp/slices"
	"net/url"
	"os"
	"strings"
)

func WriteToFile(filename string, data string) {
	writeFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	dataWriter := bufio.NewWriter(writeFile)

	if err != nil {
		log.Error(err)
	}
	dataWriter.WriteString(data)
	dataWriter.Flush()
	writeFile.Close()
}

func WriteHostsToJSONFile(filename string, data []Host) {
	file, _ := json.MarshalIndent(data, "", " ")
	WriteToFile(filename, string(file))
}

func AppendIfMissing(slice []string, key string) []string {
	for _, element := range slice {
		if element == key {
			return slice
		}
	}
	return append(slice, key)
}

func GetHost(str string) string {
	u, err := url.Parse(str)
	if err != nil {
		return str
	}
	return u.Scheme + "://" + u.Host
}

func ConvertJSONLtoJSON(input string) string {

	var data []byte
	data = append(data, '[')

	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	isFirst := true
	for _, line := range lines {
		if !isFirst && strings.TrimSpace(line) != "" {
			data = append(data, ',')
			data = append(data, '\n')
		}
		if strings.TrimSpace(line) != "" {
			data = append(data, line...)
		}
		isFirst = false
	}
	data = append(data, ']')
	return string(data)
}

func checkIfPortIsContained(port int, portSlice []int) bool {
	if slices.Contains(portSlice, port) {
		return true
	}

	return false
}
