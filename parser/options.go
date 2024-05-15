package parser

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/projectdiscovery/goflags"
	folderutil "github.com/projectdiscovery/utils/folder"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var (
	defaultSettingsLocation = filepath.Join(folderutil.HomeDirOrDefault("."), ".config/nmapParser/settings.yaml")
)

type Options struct {
	SettingsFile string
	Project      string
	ServiceJSON  bool
	HTTPDomains  bool
	All          bool
	BaseFolder   string
	Silent       bool
	Version      bool
	NoColor      bool
	Verbose      bool
}

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	options := &Options{}
	var err error
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`get simple findings from the obtained information for the specified project`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&options.Project, "project", "p", "", "project name for metadata addition"),
		flagSet.BoolVarP(&options.ServiceJSON, "serviceJSON", "sj", false, "creates a service JSON file in the project folder"),
		flagSet.BoolVarP(&options.HTTPDomains, "http_domains", "hd", false, "create additional combinations from ports for HTTP scanning"),
		flagSet.BoolVarP(&options.All, "all", "a", true, "create all combinations from ports and hosts"),
	)

	flagSet.CreateGroup("debug", "Debug",
		flagSet.BoolVar(&options.Silent, "silent", false, "show only results in output"),
		flagSet.BoolVar(&options.Version, "version", false, "show version of the project"),
		flagSet.BoolVar(&options.Verbose, "v", false, "show verbose output"),
		flagSet.BoolVarP(&options.NoColor, "no-color", "nc", false, "disable colors in output"),
	)

	flagSet.CreateGroup("config", "Config",
		flagSet.StringVar(&options.SettingsFile, "config", defaultSettingsLocation, "settings (Yaml) file location"),
	)

	if err := flagSet.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	options.configureOutput()

	if options.Version {
		fmt.Printf("Current Version: %s\n", VERSION)
		os.Exit(0)
	}

	// Validate the options passed by the user and if any
	// invalid options have been used, exit.
	err = options.validateOptions()
	if err != nil {
		log.Fatalf("Program exiting: %v\n", err)
	}

	return options
}

func (options *Options) configureOutput() {
	if options.Verbose {
		log.SetLevel(logrus.TraceLevel)
	}

	if options.NoColor {
		log.SetFormatter(&logrus.TextFormatter{
			PadLevelText:     true,
			ForceColors:      false,
			DisableTimestamp: true,
		})
	}

	if options.Silent {
		log.SetLevel(logrus.PanicLevel)
	}
}

// validateOptions validates the configuration options passed
func (options *Options) validateOptions() error {

	// Both verbose and silent flags were used
	if options.Verbose && options.Silent {
		return errors.New("both verbose and silent mode specified")
	}

	return nil
}
