package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mikemackintosh/twonicorn/config"
	yaml "gopkg.in/yaml.v3"
)

var (
	flagConfigFile string
	flagShowKeys   bool
)

func init() {
	flag.StringVar(&flagConfigFile, "c", "", "Configuration file to validate")
	flag.BoolVar(&flagShowKeys, "keys", false, "Shows config keys or not")
}

func main() {
	flag.Parse()

	if len(flagConfigFile) == 0 {
		fmt.Println("Please provide a valid configuration file")
		os.Exit(0)
	}

	r, err := ioutil.ReadFile(flagConfigFile)
	if err != nil {
		fmt.Printf("Error reading configuration file:\n\t%s\n", err)
		os.Exit(1)
	}

	var entries config.Entries
	err = yaml.Unmarshal(r, &entries)
	if err != nil {
		fmt.Printf("Error parsing configuration file:\n\t%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d configuration(s).\n\n", len(entries))

	fmt.Println("Validating configurations:")
	for name, entry := range entries {
		fmt.Printf(" - %s ", name)
		err = entry.Validate()
		if err != nil {
			fmt.Printf("[Error] - (%s)\n", err)
		} else {
			fmt.Println("[OK]")
		}

		// Show the config keys if passed
		if flagShowKeys {
			fmt.Printf("\t Key: %s\n", config.ComputeKey(name))
		}
	}
}
