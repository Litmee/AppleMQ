package apple

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// configuration item storage map
var options = make(map[string]string)

func initOptions() {
	options["mode"] = ""
	options["work-mode"] = ""
	options["port"] = ""
	options["cluster-port"] = ""
	options["cluster-map"] = ""
}

// Read the apple.ini configuration file in the conf directory
func readConf() {

	// Get the current directory of the project
	currentDir, _ := os.Getwd()
	index := strings.LastIndex(currentDir, "\\bin")

	// Assemble the path to the configuration files that AppleMQ needs ideally
	confDir := string([]byte(currentDir)[:index]) + "\\conf\\apple.ini"

	// File path validity judgment
	fi, err := os.Stat(confDir)

	// Exception and file judgment
	if err != nil || fi.IsDir() {
		panic("The apple.ini configuration file is missing from: " + confDir)
	}

	f, _ := os.Open(confDir)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("The apple.ini configuration file closed abnormally: ", err.Error())
		}
	}(f)

	// Initialize options map
	initOptions()

	// read configuration file
	buf := bufio.NewReader(f)
	for {
		str, err := buf.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				panic("An exception occurred during the content reading of the apple.ini configuration file")
			}
		}
		// remove spaces at both ends
		str = strings.TrimSpace(str)

		if str == "" || strings.HasPrefix(str, "#") {
			if err == io.EOF {
				break
			}
			continue
		}

		n := strings.Index(str, "=")

		_, ok := options[str[:n]]

		if !ok {
			panic("There is an unknown configuration item in the apple.ini configuration file, which may lead to uncertain problems after the system is running, please correct the parameters of this configuration item: " + str[:n])
		}

		options[str[:n]] = str[n+1:]
	}
}
