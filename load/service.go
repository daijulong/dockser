package load

import (
	"github.com/daijulong/dockser/lib"
	"strings"
)

func Services(services []string) (string, error) {
	if len(services) < 1 {
		return "", nil
	}
	servicesContentBlocks := make([]string, 0)
	for _, service := range services {
		file := "./docker-compose/services/" + service + ".yml"
		fileLines, err := lib.ReadFileLines(file, service)
		if err != nil {
			return "", err
		}
		outLines := make([]string, 0)
		for _, line := range fileLines {
			outLines = append(outLines, "  "+line)
		}
		servicesContentBlocks = append(servicesContentBlocks, strings.Join(outLines, "\n"))
	}
	return strings.Join(servicesContentBlocks, "\n"), nil
}
