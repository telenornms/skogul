package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/telenornms/skogul/config"
)

func getExampleConfigs() (map[string][]byte, error) {
	examplesPath := "../docs/examples"

	files, err := ioutil.ReadDir(examplesPath)

	if err != nil {
		return nil, err
	}

	bytes := make(map[string][]byte, 0)

	for _, filename := range files {
		b, err := readFileAndParseConfig(examplesPath, filename)
		if err != nil {
			return nil, err
		}

		// continue if we don't get an error but we don't get any bytes either
		// we skip some files that are not json etc.
		if b == nil {
			continue
		}

		bytes[filename.Name()] = b
	}

	return bytes, nil
}

func readFileAndParseConfig(path string, info os.FileInfo) ([]byte, error) {
	// Assuming we can parse all .json files in the example config directory
	if filepath.Ext(info.Name()) != ".json" {
		return nil, nil
	}

	data, err := ioutil.ReadFile(filepath.Join(path, info.Name()))

	if err != nil {
		return nil, err
	}

	return data, nil
}

func TestExampleConfigs(t *testing.T) {
	bytebytes, err := getExampleConfigs()

	if err != nil {
		t.Errorf("Failed to read configuration files, %s", err)
	}

	for filename, bytes := range bytebytes {

		// Log the filename in case we get a warning from the
		// configuration function so we can see what file generated the warning
		// If the test passes it gets suppressed.
		logrus.Debugf("Parsing %s", filename)
		conf, err := config.Bytes(bytes)

		if err != nil {
			t.Errorf("Failed to parse config in %s %s", filename, err)
			return
		}

		if conf == nil {
			t.Errorf("Configuration was nil for %s", filename)
			return
		}
	}
}
