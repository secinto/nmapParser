package main

import (
	"github.com/projectdiscovery/gologger"
	"github.com/secinto/nmapParser/parser"
)

func main() {
	// Parse the command line flags and read config files
	options := parser.ParseOptions()

	newParser, err := parser.NewParser(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create differ: %s\n", err)
	}

	err = newParser.Parse()
	if err != nil {
		gologger.Fatal().Msgf("Could not diff: %s\n", err)
	}
}
